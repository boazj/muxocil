// Package termprobe helps identify the terminal emulator or terminal multiplexer according to the ECMA-48 spec
// in addition to some fun huristcs.
//
// This package should be used to spoof terminals but rather to optimize behavior of tools to the working terminal
package termprobe

//lint:file-ignore ST1003

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/tiendc/gofn"
)

func SendSequence(action ProbeActions, seq []byte, w io.Writer, r io.Reader) ([]byte, error) {
	_, err := w.Write(seq)
	if err != nil {
		return nil, SendFailedError(action, err)
	}
	buf := make([]byte, 256)
	nb, err := r.Read(buf)
	if err != nil {
		return nil, ReadFailedError(action, err)
	}
	return buf[:nb], nil
}

func SendXTVERSION(w io.Writer, r io.Reader) ([]byte, error) {
	return SendSequence(XtVersion, XTVERSION, w, r)
}

func SendXTGETTCAP(w io.Writer, r io.Reader) ([]byte, error) {
	return SendSequence(XtGetTcap, XTGETTCAP, w, r)
}

func SendDA1(w io.Writer, r io.Reader) ([]byte, error) {
	return SendSequence(PDA, DA1, w, r)
}

func SendDA2(w io.Writer, r io.Reader) ([]byte, error) {
	return SendSequence(SDA, DA2, w, r)
}

func SendDA3(w io.Writer, r io.Reader) ([]byte, error) {
	return SendSequence(TDA, DA3, w, r)
}

func ParseXTVERSIONResponse(b []byte) (string, error) {
	// XTVERSION
	// DCS > | text ST
	if len(b) == 0 {
		return "", EmptyResponseError(XtVersion)
	}
	badEsc := gofn.MapSliceToMapKeys(gofn.Filter(C1, func(t byte) bool {
		return t != C1DCS && t != C1ST
	}), true)

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
			if c == '>' {
				prev = c
				continue
			}
			if c == '|' && prev == '>' {
				prefix = true
			} else if gofn.MapGet(badEsc, c, false) {
				return "", UnexpectedEscapeCodesError(XtVersion)
			} else {
				return "", MissingPrefixError(XtVersion, "DCS > |")
			}
		} else if !termination {
			if c == C1ST || (prev == ST[0] && c == ST[1]) {
				termination = true
				break
			} else if (prev == ST[0] && c != ST[1]) || gofn.MapGet(badEsc, c, false) {
				return "", UnexpectedEscapeCodesError(XtVersion)
			} else if c != ST[0] {
				output = append(output, c)
			}
		}
		prev = c
	}
	if !control || !prefix {
		return "", MissingPrefixError(XtVersion, "DCS > |")
	}
	if !termination {
		return "", MissingSuffixError(XtVersion, "ST")
	}
	return string(output[:]), nil
}

func ParseXTGETTCAPResponse(b []byte) (string, error) {
	// XTGETTCAP
	// DCS 1 + r Pt ST
	if len(b) == 0 {
		return "", EmptyResponseError(XtGetTcap)
	}
	badEsc := gofn.MapSliceToMapKeys(gofn.Filter(C1, func(t byte) bool {
		return t != C1DCS && t != C1ST
	}), true)

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
				return "", BadRequestError(XtGetTcap)
			}
		} else if !prefix {
			if c == '+' {
				prev = c
				continue
			}
			if c == 'r' && prev == '+' {
				prefix = true
			} else if gofn.MapGet(badEsc, c, false) {
				return "", UnexpectedEscapeCodesError(XtGetTcap)
			} else {
				return "", MissingPrefixError(XtGetTcap, "DCS 1 + r")
			}
		} else if !termination {
			if c == C1ST || (prev == ST[0] && c == ST[1]) {
				termination = true
				break
			} else if (prev == ST[0] && c != ST[1]) || gofn.MapGet(badEsc, c, false) {
				return "", UnexpectedEscapeCodesError(XtGetTcap)
			} else if c != ST[0] {
				output = append(output, c)
			}
		}
		prev = c
	}
	if !control || !prefix {
		return "", MissingPrefixError(XtGetTcap, "DCS 1 + r")
	}
	if !termination {
		return "", MissingSuffixError(XtGetTcap, "ST")
	}

	body := ""
	if len(output) >= 5 && output[0] == TN[0] && output[1] == TN[1] && output[2] == TN[2] && output[3] == TN[3] && output[4] == '=' {
		body = string(output[5:])
	} else {
		return "", BadResponseBodyError(XtGetTcap, "XTGETTCAP malformed terminal capability body")
	}

	return body, nil
}

func ParseDA1Response(b []byte) (string, string, error) {
	// CSI ? Tid ; Ps c
	// Tid standard VT id code
	// Ps semicolon separated parameters
	// TODO: implement
	return "", "", nil
}

// ParseDA2Response parses the response sequence of a secondary device attribute request to a terminal
// returns Pp, Pv, Pc, error
// Pp - terminal type
// Pv - firmware version (per spec, in emulators it's application version)
// Pc - ROM cartridge registration number (per spec should always be zero)
func ParseDA2Response(b []byte) (string, string, string, error) {
	// CSI  > Pp ; Pv ; Pc c
	// TODO: implement tests

	if len(b) == 0 {
		return "", "", "", EmptyResponseError(SDA)
	}
	badEsc := gofn.MapSliceToMapKeys(gofn.Filter(C1, func(t byte) bool {
		return t != C1CSI
	}), true)

	colons := 0
	control := false
	prefix := false
	termination := false

	var prev byte
	output := make([]byte, 0)

	for i := range b {
		c := b[i]
		if !control {
			if (c == C1CSI) || (c == CSI[1] && prev == CSI[0]) {
				control = true
			}
		} else if !prefix {
			if c == '>' {
				prefix = true
			} else if gofn.MapGet(badEsc, c, false) {
				return "", "", "", UnexpectedEscapeCodesError(SDA)
			} else {
				return "", "", "", MissingPrefixError(SDA, "CSI >")
			}
		} else if !termination {
			if c == ';' {
				colons++
			} else if c == 'c' && colons == 2 {
				termination = true
				break
			} else if gofn.MapGet(badEsc, c, false) {
				return "", "", "", UnexpectedEscapeCodesError(SDA)
			}
			output = append(output, c)
		}
		prev = c
	}
	if !control || !prefix {
		return "", "", "", MissingPrefixError(SDA, "CSI >")
	}
	if !termination {
		return "", "", "", MissingSuffixError(SDA, "c")
	}
	if colons != 2 {
		return "", "", "", BadResponseBodyError(
			SDA,
			fmt.Sprintf("DA2 response sequence issue, expected 3 values, got %d", colons),
		)
	}
	vals := strings.Split(string(output[:]), ":")
	return vals[0], vals[1], vals[2], nil
}

func ParseDA3Response(b []byte) (string, error) {
	// DCS ! |  D..D ST

	// TODO: implement tests

	if len(b) == 0 {
		return "", EmptyResponseError(TDA)
	}
	badEsc := gofn.MapSliceToMapKeys(gofn.Filter(C1, func(t byte) bool {
		return t != C1DCS && t != C1ST
	}), true)

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
			if c == '!' {
				prev = c
				continue
			}
			if c == '|' && prev == '!' {
				prefix = true
			} else if gofn.MapGet(badEsc, c, false) {
				return "", UnexpectedEscapeCodesError(TDA)
			} else {
				return "", MissingPrefixError(TDA, "DCS ! |")
			}
		} else if !termination {
			if c == C1ST || (prev == ST[0] && c == ST[1]) {
				termination = true
				break
			} else if (prev == ST[0] && c != ST[1]) || gofn.MapGet(badEsc, c, false) {
				return "", UnexpectedEscapeCodesError(TDA)
			} else if c != ST[0] {
				output = append(output, c)
			}
		}
		prev = c
	}
	if !control || !prefix {
		return "", MissingPrefixError(TDA, "DCS ! |")
	}
	if !termination {
		return "", MissingSuffixError(TDA, "ST")
	}
	return string(output[:]), nil
}

// TODO: IOCTL & ISATTY

func Probe() (*ProbeData, error) {
	goos := runtime.GOOS
	term, _ := os.LookupEnv("TERM")
	termProgram, _ := os.LookupEnv("TERM_PROGRAM")

	reader := bufio.NewReader(os.Stdin)

	buf, err := SendDA3(os.Stdout, reader)
	if err != nil {
		return nil, err // TODO: wrap
	}
	da3, err := ParseDA3Response(buf)
	if err != nil {
		return nil, err // TODO: wrap
	}

	buf, err = SendDA2(os.Stdout, reader)
	if err != nil {
		return nil, err // TODO: wrap
	}
	da2pp, da2pv, da2pc, err := ParseDA2Response(buf)
	if err != nil {
		return nil, err // TODO: wrap
	}

	buf, err = SendXTVERSION(os.Stdout, reader)
	if err != nil {
		return nil, err // TODO: wrap
	}
	xtversion, err := ParseXTVERSIONResponse(buf)
	if err != nil {
		return nil, err // TODO: wrap
	}

	buf, err = SendXTGETTCAP(os.Stdout, reader)
	if err != nil {
		return nil, err // TODO: wrap
	}
	xtgettcap, err := ParseXTGETTCAPResponse(buf)
	if err != nil {
		return nil, err // TODO: wrap
	}

	buf, err = SendDA1(os.Stdout, reader)
	if err != nil {
		return nil, err // TODO: wrap
	}
	da1pp, da1ps, err := ParseDA1Response(buf)
	if err != nil {
		return nil, err // TODO: wrap
	}

	return &ProbeData{
		OS:             goos,
		EnvTerm:        term,
		EnvTermProgram: termProgram,
		DA3:            da3,
		DA2Pp:          da2pp,
		DA2Pv:          da2pv,
		DA2Pc:          da2pc,
		XtermVersion:   xtversion,
		XtermGetTcap:   xtgettcap,
		DA1Pp:          da1pp,
		DA1Ps:          da1ps,
	}, nil
}
