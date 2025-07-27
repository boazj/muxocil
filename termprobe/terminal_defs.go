package termprobe

import (
	"fmt"
	"strings"

	"github.com/boazj/muxocil/utils"
)

var (
	Unknown = TerminalMD{"unknown", KindUnknown, []TermInfoReports{}, []ProbeActions{}, nil}

	// TDA response: "7E565445"
	// XTVERSION prefix: "VTE("
	// mismatch - Support for version added in 2024.
	GnomeVTE = TerminalMD{
		"vte",
		Xterm,
		[]TermInfoReports{ReportsXtermVersion, ReportsSDA},
		[]ProbeActions{TDA, XtVersion},
		probeGnomeVTE,
	}

	// TDA response: "7E484445"
	// XTVERSION prefix: "Konsole "
	// https://github.com/KDE/konsole/blob/bebbdcdb4713598a6d1497f70803a11ea21bd208/src/Vt102Emulation.cpp#L2483
	// mismatch - Support for version added in 2023.
	KdeKonsole = TerminalMD{
		"konsole",
		Xterm,
		[]TermInfoReports{ReportsXtermVersion},
		[]ProbeActions{TDA, XtVersion},
		probeKdeKonsole,
	}

	// TDA response: "7E7E5459"
	// XTVERSION prefix: "terminology "
	// mismatch - both work, depends on version.
	Terminology = TerminalMD{
		"terminology",
		Xterm,
		[]TermInfoReports{ReportsXtermVersion},
		[]ProbeActions{TDA, XtVersion},
		probeTerminology,
	}

	// TDA response: "464F4F54"
	// XTVERSION prefix: "foot("
	// mismatch - both work, depends on version.
	Foot = TerminalMD{
		"foot",
		Wayland,
		[]TermInfoReports{ReportsXtermVersion},
		[]ProbeActions{TDA, XtVersion},
		probeFoot,
	}

	// XTGETTCAP TN resposne: "mlterm"
	// XTVERSION prefix: "mlterm("
	// mismatch - both work, depends on version.
	Mlterm = TerminalMD{
		"mlterm",
		Xterm,
		[]TermInfoReports{ReportsXtermVersion},
		[]ProbeActions{XtGetTcap, XtVersion},
		probeMlterm,
	}

	// XTGETTCAP TN resposne: "xterm-kitty"
	// XTVERSION prefix: "kitty("
	// mismatch - both work, depends on version.
	Kitty = TerminalMD{
		"kitty",
		OpenGl,
		[]TermInfoReports{ReportsXtermVersion, ReportsSDA},
		[]ProbeActions{XtGetTcap, XtVersion},
		probeKitty,
	}

	// XTGETTCAP TN resposne: "xterm-ghostty"
	// XTVERSION prefix: "ghostty "
	// mismatch - both work, depends on version.
	Ghostty = TerminalMD{
		"ghostty",
		Miscellaneous,
		[]TermInfoReports{ReportsXtermVersion},
		[]ProbeActions{XtVersion, XtGetTcap},
		probeGhostty,
	}

	// SDA response: "\x1b[>0;version;1c"
	// TERM value "alacritty", can't fully trust TERM.
	Alacritty = TerminalMD{
		"alacritty",
		OpenGl,
		[]TermInfoReports{ReportsSDA},
		[]ProbeActions{SDA, EnvTerm},
		probeAlacritty,
	}

	// SDA response: "\x1b[>83;version;"
	// TODO: verfiy response structure is compliant
	// TODO: mismatch.
	GnuScreen = TerminalMD{
		"gnuscreen",
		UNIX,
		[]TermInfoReports{},
		[]ProbeActions{SDA},
		probeGnuScreen,
	}

	// XTVERSION prefix: "XTerm(".
	GeneralXterm = TerminalMD{
		"xterm",
		Xterm,
		[]TermInfoReports{ReportsXtermVersion},
		[]ProbeActions{XtVersion},
		probeXterm,
	}

	// XTVERSION prefix: "WezTerm ".
	Wezterm = TerminalMD{
		"wezterm",
		Miscellaneous,
		[]TermInfoReports{ReportsXtermVersion, ReportsSDA},
		[]ProbeActions{XtVersion},
		probeWezterm,
	}

	// XTVERSION prefix: "contour ".
	Contour = TerminalMD{
		"contour",
		Miscellaneous,
		[]TermInfoReports{ReportsXtermVersion},
		[]ProbeActions{XtVersion},
		probeContour,
	}

	// XTVERSION prefix: "tmux ".
	Tmux = TerminalMD{
		"tmux",
		UNIX,
		[]TermInfoReports{ReportsXtermVersion},
		[]ProbeActions{XtVersion},
		probeTmux,
	}

	// XTVERSION prefix: "iTerm2 "
	// Should have a separate entry for item vs iterm2?
	Iterm2 = TerminalMD{
		"iterm2",
		Apple,
		[]TermInfoReports{ReportsXtermVersion, ReportsSDA},
		[]ProbeActions{XtVersion},
		probeIterm2,
	}

	// XTVERSION prefix: "mintty "
	// NOTE: windows.
	Mintty = TerminalMD{
		"mintty",
		Microsoft,
		[]TermInfoReports{ReportsXtermVersion},
		[]ProbeActions{XtVersion, OS},
		probeMintty,
	}

	// XTVERSION prefix: "Zellij("
	// https://github.com/zellij-org/zellij/blob/48ecb0e34ff9d6d04f574237dd3a8e18e2830e6c/zellij-server/src/panes/grid.rs#L3226
	// NOTE: darwin, linux.
	Zellij = TerminalMD{
		"terminal",
		KindUnknown,
		[]TermInfoReports{},
		[]ProbeActions{XtVersion, OS},
		probeZellij,
	}

	RXVT = TerminalMD{
		"rxvt",
		Xterm,
		[]TermInfoReports{},
		[]ProbeActions{EnvTerm},
		probeRxvt,
	}

	// NOTE: darwin.
	TerminalApp = TerminalMD{
		"terminal.app",
		Apple,
		[]TermInfoReports{},
		[]ProbeActions{EnvTermProgram, OS},
		probeTerminalApp,
	}

	// TDA response: "\x1bP!|00000000\x1b\\"
	// SDA response: "\x1b[>0;10;1c"
	// Not good enough, but together it's probably ok until they will implement XTVERSION
	// NOTE: windows
	// TODO: not good enough
	// WindowsTerminal = TerminalMD{"ms-terminal", Microsoft, []TermInfoReports{ReportsSDA}, []ProbingStrategy{TDA, SDA, OS}}.

	// SDA response: "\x1b[>0;135;0c" not good enough
	// TODO: not good enough
	// Putty = TerminalMD{"putty", Microsoft, []TermInfoReports{ReportsSDA}, []ProbingStrategy{SDA}}.

	// TODO: unknown
	// DomTerm = TerminalMD{"domterm", Web, []TermInfoReports{ReportsXtermVersion, ReportsSDA}, []ProbingStrategy{}}.

	// TODO: unknown
	// TeraTerm = TerminalMD{"teraterm", Microsoft, []TermInfoReports{ReportsSDA}, []ProbingStrategy{}}.

	// TODO: unknown
	// Rlogin = TerminalMD{"rlogin", Microsoft, []TermInfoReports{ReportsXtermVersion, ReportsSDA}, []ProbingStrategy{}}.

	// TODO: unknown
	// Vscode = TerminalMD{"vscode", Microsoft, []TermInfoReports{ReportsSDA}, []ProbingStrategy{}}.
)

func probeGnomeVTE(probe *ProbeData) (bool, error) {
	if probe.DA3 == "7E565445" {
		if !strings.HasPrefix(probe.XtermVersion, "VTE(") {
			// TODO: log old VTE, but not error, Support for version added in 2024
			return true, nil
		}
		return true, nil
	}
	return false, nil
}

func probeKdeKonsole(probe *ProbeData) (bool, error) {
	if probe.DA3 == "7E484445" {
		if !strings.HasPrefix(probe.XtermVersion, "Konsole ") {
			// TODO: log old konsole, but not error, Support for version added in 2023
			return true, nil
		}
		return true, nil
	}
	return false, nil
}

func probeTerminology(probe *ProbeData) (bool, error) {
	return probe.DA3 == "7E7E5459" || strings.HasPrefix(probe.XtermVersion, "terminology "), nil
}

func probeFoot(probe *ProbeData) (bool, error) {
	return probe.DA3 == "464F4F54" || strings.HasPrefix(probe.XtermVersion, "foot("), nil
}

func probeMlterm(probe *ProbeData) (bool, error) {
	return probe.XtermGetTcap == "mlterm" || strings.HasPrefix(probe.XtermVersion, "mlterm("), nil
}

func probeGhostty(probe *ProbeData) (bool, error) {
	return probe.XtermGetTcap == "xterm-ghostty" || strings.HasPrefix(probe.XtermVersion, "ghostty "), nil
}

func probeAlacritty(probe *ProbeData) (bool, error) {
	// TODO: implement
	return false, fmt.Errorf("not implemented yet")
}

func probeGnuScreen(probe *ProbeData) (bool, error) {
	// TODO: implement
	// ver < 10000 -> NOT SCREEN
	// int s = snprintf(verstr, sizeof(verstr), "%u.%02u.%02u", ver / 10000, ver / 100 % 100, ver % 100);
	// s < 0 || (unsigned)s >= sizeof(verstr) -> NOT SCREEN
	// https://github.com/dankamongmen/notcurses/blob/94e36ccc32ed65aa394a895d53c46081aaef4450/src/lib/in.c#L1336
	return false, fmt.Errorf("not implemented yet")
}

func probeXterm(probe *ProbeData) (bool, error) {
	return strings.HasPrefix(probe.XtermVersion, "XTerm("), nil
}

func probeKitty(probe *ProbeData) (bool, error) {
	if probe.OS == utils.Windows {
		return false, nil
	}
	return probe.XtermGetTcap == "xterm-kitty" && strings.HasPrefix(probe.XtermVersion, "kitty("), nil
}

func probeTmux(probe *ProbeData) (bool, error) {
	// TODO: what about the underlying term? is there a way to identify it?
	return strings.HasPrefix(probe.XtermVersion, "tmux "), nil
}

func probeContour(probe *ProbeData) (bool, error) {
	return strings.HasPrefix(probe.XtermVersion, "contour "), nil
}

func probeWezterm(probe *ProbeData) (bool, error) {
	return strings.HasPrefix(probe.XtermVersion, "WezTerm "), nil
}

func probeIterm2(probe *ProbeData) (bool, error) {
	return probe.OS == utils.Darwin && strings.HasPrefix(probe.XtermVersion, "iTerm2 "), nil
}

func probeMintty(probe *ProbeData) (bool, error) {
	return probe.OS == utils.Windows && strings.HasPrefix(probe.XtermVersion, "mintty "), nil
}

func probeZellij(probe *ProbeData) (bool, error) {
	return probe.OS != utils.Windows && strings.HasPrefix(probe.XtermVersion, "Zellij("), nil
}

func probeRxvt(probe *ProbeData) (bool, error) {
	return strings.HasPrefix(probe.EnvTerm, "rxvt"), nil
}

func probeTerminalApp(probe *ProbeData) (bool, error) {
	return probe.OS == utils.Darwin && probe.EnvTermProgram == "Apple_Terminal", nil
}
