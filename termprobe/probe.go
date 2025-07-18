package termprobe

//lint:file-ignore ST1003

import (
	"fmt"
	"os"
	"runtime"

	"github.com/tiendc/gofn"
)

// https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h3-PC-Style-Function-Keys

// Keys
const (
	ESC byte = 0x1b
)

// C0 (7-bit) Control Characters
var (
	// Device Control String
	DCS = []byte{ESC, 'P'}

	// Control Sequence Introducer
	CSI = []byte{ESC, '['}

	// String Terminator
	ST = []byte{ESC, '\\'}
)

// C1 (8-bit) Control Characters
var (
	// Device Control String
	C1DCS byte = '\x90'

	// Control Sequence Introducer
	C1CSI byte = '\x9b'

	// String Terminator
	C1ST byte = '\x9d'
)

// Functions using CSI
var (

	// CSI > Ps q
	//             Ps = 0  ⇒  Report xterm name and version (XTVERSION).
	//           The response is a DSR sequence identifying the version:
	//             DCS > | text ST
	XTVERSION = gofn.Concat(CSI, []byte{'>', Ps, 'q'})

	// See XTVERSION, only using C1 based codes
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
	// See XTGETTCAP, only using C1 based codes
	C1XTGETTCAP = gofn.Concat([]byte{C1DCS, '+', 'q'}, TN, []byte{C1ST})
	TN          = []byte{'5', '4', '4', 'e'}

	// CSI = Ps c
	//         Send Device Attributes (Tertiary DA).
	//           Ps = 0  ⇒  report Terminal Unit ID (default), VT400.  XTerm
	//         uses zeros for the site code and serial number in its DECRPTUI
	//         response.
	DA3 = gofn.Concat(CSI, []byte{'=', Ps, 'c'})
	// See DA3, only using C1 based codes
	C1DA3 = []byte{C1CSI, '=', Ps, 'c'}

	// CSI > Ps c
	//         Send Device Attributes (Secondary DA).
	//           Ps = 0  or omitted ⇒  request the terminal's identification
	//         code.  The response depends on the decTerminalID resource
	//         setting.  It should apply only to VT220 and up, but xterm
	//         extends this to VT100.
	DA2 = gofn.Concat(CSI, []byte{'>', Ps, 'c'})
	// See DA2, only using C1 based codes
	C1DA2 = []byte{C1CSI, '>', Ps, 'c'}

	// CSI Ps c  Send Device Attributes (Primary DA).
	//           Ps = 0  or omitted ⇒  request attributes from terminal.  The
	//         response depends on the decTerminalID resource setting.
	DA1 = gofn.Concat(CSI, []byte{Ps, 'c'})
	// See DA1, only using C1 based codes
	C1DA1 = []byte{C1CSI, Ps, 'c'}

	// Common attribute used in multiple control sequences used to get information
	// form the terminal device
	//
	// PDA - request attributes from terminal
	// SDA - request the terminal's identification code
	// TDA - report Terminal Unit ID
	// XTVERSION - Report xterm name and version
	Ps byte = '0'
)

func ParseXTVERSIONResponse(b []byte) (string, error) {
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
			if (b[i] == C1DCS) || (b[i] == DCS[1] && prev == DCS[0]) {
				control = true
			}
		} else if !prefix {
			if b[i] == '|' && prev == '>' {
				prefix = true
			}
		} else if !termination {
			if b[i] == C1ST || (prev == ST[0] && b[i] == ST[1]) {
				termination = true
				break
			} else if prev == ST[0] && b[i] != ST[1] {
				// unexpected escape control
				break
			} else if b[i] != ST[0] {
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

func ParseXTGETTCAPResponse(b []byte) (string, error) {
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
			if (b[i] == C1DCS) || (b[i] == DCS[1] && prev == DCS[0]) {
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
			if b[i] == C1ST || (prev == ST[0] && b[i] == ST[1]) {
				termination = true
				break
			} else if prev == ST[0] && b[i] != ST[1] {
				// unexpected escape control
				break
			} else if b[i] != ST[0] {
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

type TerminalKind int

// Definition according to terminfo db https://invisible-island.net/ncurses/terminfo.src-sections.htm
// NOTE: additional terms defined
// Xterm: Kterm ETERM ATERM XITERM GPTERM EMU MVTERM MTERM VWM MGR SimpleTerm TERMINATOR
// Miscellaneous: pangoterm
// UNIX: Mosh Dvtm Screen Emacs
// NonUNIX Consoles: Cygwin
const (
	Xterm TerminalKind = iota
	OpenGl
	Wayland
	UNIX
	NonUNIX
	Web
	Apple
	Microsoft
	Miscellaneous
	KindUnknown
)

type TermInfoReports int

// Definition according to terminfo db https://invisible-island.net/ncurses/terminfo.src-entries.html
const (
	ReportsXtermVersion TermInfoReports = iota
	ReportsSDA
)

type ProbingStrategy int

// Definition according to notcurses https://github.com/dankamongmen/notcurses/blob/master/src/lib/in.h#L31
// Others have been added by specific testing
const (
	None ProbingStrategy = iota
	XtVersion
	XtGetTcap
	PDA
	SDA
	TDA
	EnvTerm
	EnvTermProgram
	OS
)

type TerminalMD struct {
	ID       string
	Kind     TerminalKind      // According to Terminfo
	reports  []TermInfoReports // According to Terminfo
	strategy []ProbingStrategy // According to notcurses & testing
}

var (
	Unknown = TerminalMD{"unknown", KindUnknown, []TermInfoReports{}, []ProbingStrategy{}}

	// TDA response: "\x1bP!|7E565445\x1b\\"   -> DCS7 ! | DDDDDDDD (4 hex pairs) ST7
	// XTVERSION prefix: "VTE("
	// mismatch - Support for version added in 2024
	GnomeVTE = TerminalMD{"vte", Xterm, []TermInfoReports{ReportsXtermVersion, ReportsSDA}, []ProbingStrategy{TDA, XtVersion}}

	// TDA response: "\x1bP!|7E484445\x1b\\"   -> DCS7 ! | DDDDDDDD (4 hex pairs) ST7
	// XTVERSION prefix: "Konsole "
	// https://github.com/KDE/konsole/blob/bebbdcdb4713598a6d1497f70803a11ea21bd208/src/Vt102Emulation.cpp#L2483
	// mismatch - Support for version added in 2023
	KdeKonsole = TerminalMD{"konsole", Xterm, []TermInfoReports{ReportsXtermVersion}, []ProbingStrategy{TDA, XtVersion}}

	// TDA response: "\x1bP!|7E7E5459\x1b\\"   -> DCS7 ! | DDDDDDDD (4 hex pairs) ST7
	// XTVERSION prefix: "terminology "
	// mismatch - both work, depends on version
	Terminology = TerminalMD{"terminology", Xterm, []TermInfoReports{ReportsXtermVersion}, []ProbingStrategy{TDA, XtVersion}}

	// TDA response: "\x1bP!|464F4F54\x1b\\"   -> DCS7 ! | DDDDDDDD (4 hex pairs) ST7
	// XTVERSION prefix: "foot("
	// mismatch - both work, depends on version
	Foot = TerminalMD{"foot", Wayland, []TermInfoReports{ReportsXtermVersion}, []ProbingStrategy{TDA, XtVersion}}

	// XTGETTCAP TN resposne "mlterm"
	// XTVERSION prefix: "mlterm("
	// mismatch - both work, depends on version
	Mlterm = TerminalMD{"mlterm", Xterm, []TermInfoReports{ReportsXtermVersion}, []ProbingStrategy{XtGetTcap, XtVersion}}

	// XTGETTCAP TN resposne "xterm-kitty"
	// XTVERSION prefix: "kitty("
	// mismatch - both work, depends on version
	Kitty = TerminalMD{"kitty", OpenGl, []TermInfoReports{ReportsXtermVersion, ReportsSDA}, []ProbingStrategy{XtGetTcap, XtVersion}}

	// XTGETTCAP TN resposne "xterm-ghostty"
	// XTVERSION prefix: "ghostty "
	// mismatch - both work, depends on version
	Ghostty = TerminalMD{"ghostty", Miscellaneous, []TermInfoReports{ReportsXtermVersion}, []ProbingStrategy{XtVersion, XtGetTcap}}

	// SDA response: "\x1b[>0;version;1c" -> "CSI7 > Pp(terminal type, ignore) ; Pv(value) ; Pc(ignore) c". TERM value "alacritty", can't fully trust TERM
	Alacritty = TerminalMD{"alacritty", OpenGl, []TermInfoReports{ReportsSDA}, []ProbingStrategy{SDA, EnvTerm}}

	// SDA response: "\x1b[>83;version;"  FIXME: missing data in response
	// ver < 10000 -> NOT SCREEN
	// int s = snprintf(verstr, sizeof(verstr), "%u.%02u.%02u", ver / 10000, ver / 100 % 100, ver % 100);
	// s < 0 || (unsigned)s >= sizeof(verstr) -> NOT SCREEN
	// https://github.com/dankamongmen/notcurses/blob/94e36ccc32ed65aa394a895d53c46081aaef4450/src/lib/in.c#L1336
	// TODO: mismatch
	GnuScreen = TerminalMD{"gnuscreen", UNIX, []TermInfoReports{}, []ProbingStrategy{SDA}}

	// XTVERSION prefix: "XTerm("
	GeneralXterm = TerminalMD{"xterm", Xterm, []TermInfoReports{ReportsXtermVersion}, []ProbingStrategy{XtVersion}}

	// XTVERSION prefix: "WezTerm "
	Wezterm = TerminalMD{"wezterm", Miscellaneous, []TermInfoReports{ReportsXtermVersion, ReportsSDA}, []ProbingStrategy{XtVersion}}

	// XTVERSION prefix: "contour "
	Contour = TerminalMD{"contour", Miscellaneous, []TermInfoReports{ReportsXtermVersion}, []ProbingStrategy{XtVersion}}

	// XTVERSION prefix: "tmux "
	Tmux = TerminalMD{"tmux", UNIX, []TermInfoReports{ReportsXtermVersion}, []ProbingStrategy{XtVersion}}

	// XTVERSION prefix: "iTerm2 "
	// Should have a separate entry for item vs iterm2?
	Iterm2 = TerminalMD{"iterm2", Apple, []TermInfoReports{ReportsXtermVersion, ReportsSDA}, []ProbingStrategy{XtVersion}}

	// XTVERSION prefix: "mintty "
	// NOTE: windows
	Mintty = TerminalMD{"mintty", Microsoft, []TermInfoReports{ReportsXtermVersion}, []ProbingStrategy{XtVersion, OS}}

	// XTVERSION prefix: "Zellij("
	// https://github.com/zellij-org/zellij/blob/48ecb0e34ff9d6d04f574237dd3a8e18e2830e6c/zellij-server/src/panes/grid.rs#L3226
	// NOTE: darwin, linux
	Zellij = TerminalMD{"terminal", KindUnknown, []TermInfoReports{}, []ProbingStrategy{XtVersion, OS}}

	RXVT = TerminalMD{"rxvt", Xterm, []TermInfoReports{}, []ProbingStrategy{EnvTerm}}

	// NOTE: darwin
	TerminalApp = TerminalMD{"terminal.app", Apple, []TermInfoReports{}, []ProbingStrategy{EnvTermProgram, OS}}

	// TDA response: "\x1bP!|00000000\x1b\\"
	// SDA response: "\x1b[>0;10;1c"
	// Not good enough, but together it's probably ok until they will implement XTVERSION
	// NOTE: windows
	// TODO: not good enough
	WindowsTerminal = TerminalMD{"ms-terminal", Microsoft, []TermInfoReports{ReportsSDA}, []ProbingStrategy{TDA, SDA, OS}}

	// SDA response: "\x1b[>0;135;0c" not good enough
	// TODO: not good enough
	Putty = TerminalMD{"putty", Microsoft, []TermInfoReports{ReportsSDA}, []ProbingStrategy{SDA}}

	// TODO: unknown
	DomTerm = TerminalMD{"domterm", Web, []TermInfoReports{ReportsXtermVersion, ReportsSDA}, []ProbingStrategy{}}

	// TODO: unknown
	TeraTerm = TerminalMD{"teraterm", Microsoft, []TermInfoReports{ReportsSDA}, []ProbingStrategy{}}

	// TODO: unknown
	Rlogin = TerminalMD{"rlogin", Microsoft, []TermInfoReports{ReportsXtermVersion, ReportsSDA}, []ProbingStrategy{}}

	// TODO: unknown
	Vscode = TerminalMD{"vscode", Microsoft, []TermInfoReports{ReportsSDA}, []ProbingStrategy{}}
)

// TODO: IOCTL & ISATTY

type ProbeData struct {
	OS             string
	XtermVersion   string
	XtermGetTcap   string
	PDA            string
	SDA            string
	TDA            string
	EnvTerm        string
	EnvTermProgram string
}

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
