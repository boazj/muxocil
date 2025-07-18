package termprobe

//lint:file-ignore ST1003

import (
	"fmt"

	"github.com/tiendc/gofn"
)

// Send Tertiary Device Attributes (CSI = c)
//     Identifies VTE and foot
// Send Secondary Device Attributes (CSI > c)
//     Identifies Alacritty's version number
// XTVERSION (CSI > 0 q)
//     Identifies XTerm, WezTerm, and Contour
// XTGETTCAP for the TN key (DCS + q 544e ST)
//     Identifies Kitty and MLterm
// Send Primary Device Attributes (CSI c)

// https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-PC-Style-Function-Keys

// Keys
const ()

// C0 (7-bit) Control Characters
var (
	ESC byte = 0x1b

	// Device Control String
	DCS7 = []byte{ESC, 'P'}

	// Control Sequence Introducer
	CSI7 = []byte{ESC, '['}

	// String Terminator
	ST7 = []byte{ESC, '\\'}
)

// C1 (8-bit) Control Characters
var (
	// Device Control String
	DCS8 byte = '\x90'

	// Control Sequence Introducer
	CSI8 byte = '\x9b'

	// String Terminator
	ST8 byte = '\x9d'
)

// Functions using CSI
var (

	// CSI > Ps q
	//             Ps = 0  ⇒  Report xterm name and version (XTVERSION).
	//           The response is a DSR sequence identifying the version:
	//             DCS > | text ST
	XTVERSION8 = []byte{CSI8, '>', Ps, 'q'}
	XTVERSION7 = gofn.Concat(CSI7, []byte{'>', Ps, 'q'})

	// DCS + q Pt ST
	//           Request Termcap/Terminfo String (XTGETTCAP), xterm.  The
	//           string following the "q" is a list of names encoded in
	//           hexadecimal (2 digits per character) separated by ; which
	//           correspond to termcap or terminfo capability names for special
	//           keyboard keys.  A terminal description will include other
	//           capabilities, e.g., for cursor movement, which are
	//           intentionally not part of this interface.
	//           A few more terminal capabilities are recognized, which are not
	//           names of special keys:
	//
	//           o   Co for termcap colors (or colors for terminfo colors), and
	//
	//           o   TN for termcap name (or name for terminfo name).
	//
	//           o   RGB for the ncurses direct-color extension.
	//               Only a terminfo name is provided, since termcap
	//               applications cannot use this information.
	//
	//           These capabilities fall into two categories:
	//
	//           o   Terminal capabilities which may be dynamically adjusted in
	//               xterm so they do not necessarily match a terminal
	//               description.
	//
	//           o   The name of the terminal description, which an application
	//               can use to obtain the static set of capabilities.
	//
	//           xterm responds with
	//           DCS 1 + r Pt ST for valid requests, adding to Pt an = , and
	//           the value of the corresponding string that xterm would send,
	//           or
	//           DCS 0 + r ST for invalid requests.
	//           The strings are encoded in hexadecimal (2 digits per
	//           character).  If more than one name is given, xterm replies
	//           with each name/value pair in the same response.  An invalid
	//           name (one not found in xterm's tables) ends processing of the
	//           list of names.
	XTGETTCAP8 = gofn.Concat([]byte{DCS8, '+', 'q'}, TN, []byte{ST8})
	XTGETTCAP7 = gofn.Concat(DCS7, []byte{'+', 'q'}, TN, ST7)
	TN         = []byte{'5', '4', '4', 'e'}

	// CSI = Ps c
	//         Send Device Attributes (Tertiary DA).
	//           Ps = 0  ⇒  report Terminal Unit ID (default), VT400.  XTerm
	//         uses zeros for the site code and serial number in its DECRPTUI
	//         response.
	TertiaryDA8 = []byte{CSI8, '=', Ps, 'c'}
	TertiaryDA7 = gofn.Concat(CSI7, []byte{'=', Ps, 'c'})

	// CSI > Ps c
	//         Send Device Attributes (Secondary DA).
	//           Ps = 0  or omitted ⇒  request the terminal's identification
	//         code.  The response depends on the decTerminalID resource
	//         setting.  It should apply only to VT220 and up, but xterm
	//         extends this to VT100.
	SecondaryDA8 = []byte{CSI8, '>', Ps, 'c'}
	SecondaryDA7 = gofn.Concat(CSI7, []byte{'>', Ps, 'c'})

	// CSI Ps c  Send Device Attributes (Primary DA).
	//           Ps = 0  or omitted ⇒  request attributes from terminal.  The
	//         response depends on the decTerminalID resource setting.
	PrimaryDA8 = []byte{CSI8, Ps, 'c'}
	PrimaryDA7 = gofn.Concat(CSI7, []byte{Ps, 'c'})

	// Common attribute used in multiple control sequences used to get information
	// form the terminal device
	//
	// PDA - request attributes from terminal
	// SDA - request the terminal's identification code
	// TDA - report Terminal Unit ID
	// XTVERSION - Report xterm name and version
	Ps byte = '0'
)

func XTVersionResponse(b []byte) (string, error) {
	// XTVERSION
	// DCS > | text ST
	if len(b) == 0 {
		return "", fmt.Errorf("XTVERSION response sequence is empty")
	}

	control := false
	prefix := false
	termination := false

	var prev byte
	output := make([]byte, 0)

	for i := range b {
		if !control {
			if (b[i] == DCS8) || (b[i] == DCS7[1] && prev == DCS7[0]) {
				control = true
			}
		} else if !prefix {
			if b[i] == '|' && prev == '>' {
				prefix = true
			}
		} else if !termination {
			if b[i] == ST8 || (prev == ST7[0] && b[i] == ST7[1]) {
				termination = true
				break
			} else if prev == ST7[0] && b[i] != ST7[1] {
				// unexpected escape control
				break
			} else if b[i] != ST7[0] {
				output = append(output, b[i])
			}
		}
		prev = b[i]
	}
	if !control || !prefix {
		return "", fmt.Errorf("XTVERSION response sequence DCS prefix missing")
	}
	if !termination {
		return "", fmt.Errorf("XTVERSION response sequence ST termination suffix missing")
	}
	return string(output[:]), nil
}

func XTGetTcapResponse(b []byte) (string, error) {
	// XTGETTCAP
	// DCS 1 + r Pt ST
	if len(b) == 0 {
		return "", fmt.Errorf("XTGETTCAP response sequence is empty")
	}

	control := false
	status := false
	prefix := false
	termination := false

	var prev byte
	output := make([]byte, 0)

	for i := range b {
		if !control {
			if (b[i] == DCS8) || (b[i] == DCS7[1] && prev == DCS7[0]) {
				control = true
			}
		} else if !status {
			if b[i] == '1' {
				status = true
			} else {
				// response for an illegal request
				break
			}
		} else if !prefix {
			if b[i] == 'r' && prev == '+' {
				prefix = true
			}
		} else if !termination {
			if b[i] == ST8 || (prev == ST7[0] && b[i] == ST7[1]) {
				termination = true
				break
			} else if prev == ST7[0] && b[i] != ST7[1] {
				// unexpected escape control
				break
			} else if b[i] != ST7[0] {
				output = append(output, b[i])
			}
		}
		prev = b[i]
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
