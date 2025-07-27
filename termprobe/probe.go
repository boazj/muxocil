// Package termprobe helps identify the terminal emulator or terminal multiplexer according to the ECMA-48 spec
// in addition to some fun huristcs.
//
// # This package should be used to spoof terminals but rather to optimize behavior of tools to the working terminal
package termprobe

//FIXME:
//lint:file-ignore ST1003

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/boazj/muxocil/utils"
	"github.com/tiendc/gofn"
)

var (
	xtversionBadEsc = gofn.MapSliceToMapKeys(gofn.Filter(C1, func(t byte) bool {
		return t != C1DCS && t != C1ST
	}), true)

	xtgettcapBadEsc = gofn.MapSliceToMapKeys(gofn.Filter(C1, func(t byte) bool {
		return t != C1DCS && t != C1ST
	}), true)

	da2BadEsc = gofn.MapSliceToMapKeys(gofn.Filter(C1, func(t byte) bool {
		return t != C1CSI
	}), true)

	da3BadEsc = gofn.MapSliceToMapKeys(gofn.Filter(C1, func(t byte) bool {
		return t != C1DCS && t != C1ST
	}), true)
)

func sendSequence(action ProbeActions, seq []byte, w io.Writer, r io.Reader) ([]byte, error) {
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

func sendXTVERSION(w io.Writer, r io.Reader) ([]byte, error) {
	return sendSequence(XtVersion, XTVERSION, w, r)
}

func sendXTGETTCAP(w io.Writer, r io.Reader) ([]byte, error) {
	return sendSequence(XtGetTcap, XTGETTCAP, w, r)
}

func sendDA1(w io.Writer, r io.Reader) ([]byte, error) {
	return sendSequence(PDA, DA1, w, r)
}

func sendDA2(w io.Writer, r io.Reader) ([]byte, error) {
	return sendSequence(SDA, DA2, w, r)
}

func sendDA3(w io.Writer, r io.Reader) ([]byte, error) {
	return sendSequence(TDA, DA3, w, r)
}

// nolint: dupl
func parseXTVERSIONResponse(b []byte) (string, error) {
	// XTVERSION
	// DCS > | text ST
	if len(b) == 0 {
		return "", EmptyResponseError(XtVersion)
	}

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
			} else if gofn.MapGet(xtversionBadEsc, c, false) {
				return "", UnexpectedEscapeCodesError(XtVersion)
			} else {
				return "", MissingPrefixError(XtVersion, "DCS > |")
			}
		} else if !termination {
			if c == C1ST || (prev == ST[0] && c == ST[1]) {
				termination = true
				break
			} else if (prev == ST[0] && c != ST[1]) || gofn.MapGet(xtversionBadEsc, c, false) {
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
	return string(output), nil
}

// nolint: dupl
func parseXTGETTCAPResponse(b []byte) (string, error) {
	// XTGETTCAP
	// DCS 1 + r Pt ST
	if len(b) == 0 {
		return "", EmptyResponseError(XtGetTcap)
	}

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
			} else if gofn.MapGet(xtgettcapBadEsc, c, false) {
				return "", UnexpectedEscapeCodesError(XtGetTcap)
			} else {
				return "", MissingPrefixError(XtGetTcap, "DCS 1 + r")
			}
		} else if !termination {
			if c == C1ST || (prev == ST[0] && c == ST[1]) {
				termination = true
				break
			} else if (prev == ST[0] && c != ST[1]) || gofn.MapGet(xtgettcapBadEsc, c, false) {
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

// nolint: dupl
func parseDA1Response(b []byte) (string, string, error) {
	// CSI ? Tid ; Ps c
	// Tid standard VT id code
	// Ps semicolon separated parameters
	// TODO: implement
	return "", "", nil
}

// parseDA2Response parses the response sequence of a secondary device attribute request to a terminal
// returns Pp, Pv, Pc, error
// Pp - terminal type
// Pv - firmware version (per spec, in emulators it's application version)
// Pc - ROM cartridge registration number (per spec should always be zero).
// nolint: dupl
func parseDA2Response(b []byte) (string, string, string, error) {
	// CSI  > Pp ; Pv ; Pc c
	// TODO: implement tests

	if len(b) == 0 {
		return "", "", "", EmptyResponseError(SDA)
	}

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
			} else if gofn.MapGet(da2BadEsc, c, false) {
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
			} else if gofn.MapGet(da2BadEsc, c, false) {
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
	vals := strings.Split(string(output), ":")
	return vals[0], vals[1], vals[2], nil
}

// nolint: dupl
func parseDA3Response(b []byte) (string, error) {
	// DCS ! |  D..D ST
	// TODO: implement tests
	if len(b) == 0 {
		return "", EmptyResponseError(TDA)
	}

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
			} else if gofn.MapGet(da3BadEsc, c, false) {
				return "", UnexpectedEscapeCodesError(TDA)
			} else {
				return "", MissingPrefixError(TDA, "DCS ! |")
			}
		} else if !termination {
			if c == C1ST || (prev == ST[0] && c == ST[1]) {
				termination = true
				break
			} else if (prev == ST[0] && c != ST[1]) || gofn.MapGet(da3BadEsc, c, false) {
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
	return string(output), nil
}

// TODO: IOCTL & ISATTY

func Probe() (*ProbeData, error) {
	goos := utils.GOOS()
	term, _ := os.LookupEnv("TERM")
	termProgram, _ := os.LookupEnv("TERM_PROGRAM")

	in := bufio.NewReader(os.Stdin)
	out := os.Stdout

	buf, err := sendDA3(out, in)
	if err != nil {
		return nil, err
	}
	da3, err := parseDA3Response(buf)
	if err != nil {
		if perr := IsProbeErrorOrUnknown(err, TDA); !perr.IsEmptyResponse() {
			return nil, perr
		}
	}

	buf, err = sendDA2(out, in)
	if err != nil {
		return nil, err
	}
	da2pp, da2pv, da2pc, err := parseDA2Response(buf)
	if err != nil {
		if perr := IsProbeErrorOrUnknown(err, TDA); !perr.IsEmptyResponse() {
			return nil, perr
		}
	}

	buf, err = sendXTVERSION(out, in)
	if err != nil {
		return nil, err
	}
	xtversion, err := parseXTVERSIONResponse(buf)
	if err != nil {
		if perr := IsProbeErrorOrUnknown(err, TDA); !perr.IsEmptyResponse() {
			return nil, perr
		}
	}

	buf, err = sendXTGETTCAP(out, in)
	if err != nil {
		return nil, err
	}
	xtgettcap, err := parseXTGETTCAPResponse(buf)
	if err != nil {
		if perr := IsProbeErrorOrUnknown(err, TDA); !perr.IsEmptyResponse() {
			return nil, perr
		}
	}

	buf, err = sendDA1(out, in)
	if err != nil {
		return nil, err
	}
	da1pp, da1ps, err := parseDA1Response(buf)
	if err != nil {
		if perr := IsProbeErrorOrUnknown(err, TDA); !perr.IsEmptyResponse() {
			return nil, perr
		}
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
