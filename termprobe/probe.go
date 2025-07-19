// Package termprobe helps identify the terminal emulator or terminal multiplexer according to the ECMA-48 spec
// in addition to some fun huristcs.
//
// This package should be used to spoof terminals but rather to optimize behavior of tools to the working terminal
package termprobe

//lint:file-ignore ST1003

import (
	"fmt"
	"os"
	"runtime"

	"github.com/tiendc/gofn"
)

func ParseXTVERSIONResponse(b []byte) (string, error) {
	// XTVERSION
	// DCS > | text ST
	if len(b) == 0 {
		return "", fmt.Errorf("XTVERSION response sequence is empty")
	}
	badEsc := gofn.MapSliceToMapKeys(gofn.Filter(C1, func(t byte) bool {
		return t != C1DCS && t != C1ST
	}), true)
	fmt.Println(badEsc)

	control := false
	prefix := false
	termination := false

	var prev byte
	output := make([]byte, 0)

	for i := range b {
		c := b[i]
		if !control {
			if (c == C1DCS) || (c == DCS[1] && prev == DCS[0]) {
				control = true
			}
		} else if !prefix {
			if c == '|' && prev == '>' {
				prefix = true
			} else if gofn.MapGet(badEsc, c, false) {
				// unexpected escape control
				return "", fmt.Errorf("XTVERSION response sequence contains unexpected escape code")
			}
		} else if !termination {
			if c == C1ST || (prev == ST[0] && c == ST[1]) {
				termination = true
				break
			} else if (prev == ST[0] && c != ST[1]) || gofn.MapGet(badEsc, c, false) {
				// unexpected escape control
				return "", fmt.Errorf("XTVERSION response sequence contains unexpected escape code")
			} else if c != ST[0] {
				output = append(output, c)
			}
		}
		prev = c
	}
	if !control || !prefix {
		return "", fmt.Errorf("XTVERSION response sequence DCS prefix missing")
	}
	if !termination {
		return "", fmt.Errorf("XTVERSION response sequence ST termination suffix missing")
	}
	return string(output[:]), nil
}

func ParseXTGETTCAPResponse(b []byte) (string, error) {
	// XTGETTCAP
	// DCS 1 + r Pt ST
	if len(b) == 0 {
		return "", fmt.Errorf("XTGETTCAP response sequence is empty")
	}
	badEsc := gofn.MapSliceToMapKeys(gofn.Filter(C1, func(t byte) bool {
		return t != C1DCS && t != C1ST
	}), true)
	fmt.Println(badEsc)

	control := false
	status := false
	prefix := false
	termination := false

	var prev byte
	output := make([]byte, 0)

	for i := range b {
		c := b[i]
		if !control {
			if (c == C1DCS) || (c == DCS[1] && prev == DCS[0]) {
				control = true
			}
		} else if !status {
			if c == '1' {
				status = true
			} else {
				// response for an illegal request
				break
			}
		} else if !prefix {
			if c == 'r' && prev == '+' {
				prefix = true
			} else if gofn.MapGet(badEsc, c, false) {
				// unexpected escape control
				return "", fmt.Errorf("XTGETTCAP response sequence contains unexpected escape code")
			}
		} else if !termination {
			if c == C1ST || (prev == ST[0] && c == ST[1]) {
				termination = true
				break
			} else if (prev == ST[0] && c != ST[1]) || gofn.MapGet(badEsc, c, false) {
				// unexpected escape control
				return "", fmt.Errorf("XTGETTCAP response sequence contains unexpected escape code")
			} else if c != ST[0] {
				output = append(output, c)
			}
		}
		prev = c
	}
	if control && !status {
		return "", fmt.Errorf("XTGETTCAP illegal request")
	}
	if !control || !prefix {
		return "", fmt.Errorf("XTGETTCAP response sequence DCS prefix missing")
	}
	if !termination {
		return "", fmt.Errorf("XTGETTCAP response sequence ST termination suffix missing")
	}

	body := ""
	if len(output) >= 5 && output[0] == TN[0] && output[1] == TN[1] && output[2] == TN[2] && output[3] == TN[3] && output[4] == '=' {
		body = string(output[5:])
	} else {
		return "", fmt.Errorf("XTGETTCAP malformed terminal capability body")
	}

	return body, nil
}

func ParseDA1Response(b []byte) (string, error) {
	// TODO: implement
	return "", nil
}

func ParseDA2Response(b []byte) (string, error) {
	// TODO: implement
	return "", nil
}

func ParseDA3Response(b []byte) (string, error) {
	// TODO: implement
	return "", nil
}

// TODO: IOCTL & ISATTY

func Probe() (*ProbeData, error) {
	// Start from the easy stuff
	goos := runtime.GOOS
	term, _ := os.LookupEnv("TERM")
	termProgram, _ := os.LookupEnv("TERM_PROGRAM")

	// notcurses flow
	// Send Tertiary Device Attributes (CSI = c)
	//     Identifies VTE and foot
	// Send Secondary Device Attributes (CSI > c)
	//     Identifies Alacritty's version number
	// XTVERSION (CSI > 0 q)
	//     Identifies XTerm, WezTerm, and Contour
	// XTGETTCAP for the TN key (DCS + q 544e ST)
	//     Identifies Kitty and MLterm
	// Send Primary Device Attributes (CSI c)

	return &ProbeData{
		OS:             goos,
		EnvTerm:        term,
		EnvTermProgram: termProgram,
	}, nil
}
