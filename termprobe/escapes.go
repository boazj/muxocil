package termprobe

import "github.com/tiendc/gofn"

// https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-PC-Style-Function-Keys

// Keys.
const (
	ESC byte = 0x1b
)

// C0 (7-bit) Control Characters.
var (
	// Device Control String.
	DCS = []byte{ESC, 'P'}

	// Control Sequence Introducer.
	CSI = []byte{ESC, '['}

	// String Terminator.
	ST = []byte{ESC, '\\'}
)

// C1 (8-bit) Control Characters.
var (

	// Ignore.
	C1IND byte = '\x84'

	// Ignore.
	C1NEL byte = '\x85'

	// Ignore.
	C1HTS byte = '\x88'

	// Ignore.
	C1RI byte = '\x8d'

	// Ignore.
	C1SS2 byte = '\x8e'

	// Ignore.
	C1SS3 byte = '\x8f'

	// Device Control String.
	C1DCS byte = '\x90'

	// Ignore.
	C1SPA byte = '\x96'

	// Ignore.
	C1EPA byte = '\x97'

	// Ignore.
	C1SOS byte = '\x98'

	// Ignore.
	C1DECID byte = '\x9a'

	// Control Sequence Introducer.
	C1CSI byte = '\x9b'

	// String Terminator.
	C1ST byte = '\x9c'

	// Ignore.
	C1OSC byte = '\x9d'

	// Ignore.
	C1PM byte = '\x9e'

	// Ignore.
	C1APC byte = '\x9f'

	C1 = []byte{C1IND, C1NEL, C1HTS, C1RI, C1SS2, C1SS3, C1DCS, C1SPA, C1EPA, C1SOS, C1DECID, C1CSI, C1ST, C1OSC, C1PM, C1APC}
)

// Functions using CSI.
var (

	// CSI > Ps q
	//             Ps = 0  ⇒  Report xterm name and version (XTVERSION).
	//           The response is a DSR sequence identifying the version:
	//             DCS > | text ST
	XTVERSION = gofn.Concat(CSI, []byte{'>', Ps, 'q'})

	// See XTVERSION, only using C1 based codes.
	C1XTVERSION = []byte{C1CSI, '>', Ps, 'q'}

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
	XTGETTCAP = gofn.Concat(DCS, []byte{'+', 'q'}, TN, ST)
	// See XTGETTCAP, only using C1 based codes.
	C1XTGETTCAP = gofn.Concat([]byte{C1DCS, '+', 'q'}, TN, []byte{C1ST})
	TN          = []byte{'5', '4', '4', 'e'}

	// CSI = Ps c
	//         Send Device Attributes (Tertiary DA).
	//           Ps = 0  ⇒  report Terminal Unit ID (default), VT400.  XTerm
	//         uses zeros for the site code and serial number in its DECRPTUI
	//         response.
	DA3 = gofn.Concat(CSI, []byte{'=', Ps, 'c'})
	// See DA3, only using C1 based codes.
	C1DA3 = []byte{C1CSI, '=', Ps, 'c'}

	// CSI > Ps c
	//         Send Device Attributes (Secondary DA).
	//           Ps = 0  or omitted ⇒  request the terminal's identification
	//         code.  The response depends on the decTerminalID resource
	//         setting.  It should apply only to VT220 and up, but xterm
	//         extends this to VT100.
	DA2 = gofn.Concat(CSI, []byte{'>', Ps, 'c'})
	// See DA2, only using C1 based codes.
	C1DA2 = []byte{C1CSI, '>', Ps, 'c'}

	// CSI Ps c  Send Device Attributes (Primary DA).
	//           Ps = 0  or omitted ⇒  request attributes from terminal.  The
	//         response depends on the decTerminalID resource setting.
	DA1 = gofn.Concat(CSI, []byte{Ps, 'c'})
	// See DA1, only using C1 based codes.
	C1DA1 = []byte{C1CSI, Ps, 'c'}

	// Common attribute used in multiple control sequences used to get information
	// form the terminal device
	//
	// PDA - request attributes from terminal
	// SDA - request the terminal's identification code
	// TDA - report Terminal Unit ID
	// XTVERSION - Report xterm name and version.
	Ps byte = '0'
)
