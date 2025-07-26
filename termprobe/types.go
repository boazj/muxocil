package termprobe

type TerminalKind int

// Definition according to terminfo db https://invisible-island.net/ncurses/terminfo.src-sections.htm
// NOTE: additional terms defined
// Xterm: Kterm ETERM ATERM XITERM GPTERM EMU MVTERM MTERM VWM MGR SimpleTerm TERMINATOR
// Miscellaneous: pangoterm
// UNIX: Mosh Dvtm Screen Emacs
// NonUNIX Consoles: Cygwin.
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

type ProbeActions string

// Definition according to notcurses https://github.com/dankamongmen/notcurses/blob/master/src/lib/in.h#L31
// Others have been added by specific testing.
const (
	None           ProbeActions = "None"
	XtVersion      ProbeActions = "XTVERSION"
	XtGetTcap      ProbeActions = "XTGETTCAP"
	PDA            ProbeActions = "DA1"
	SDA            ProbeActions = "DA2"
	TDA            ProbeActions = "DA3"
	EnvTerm        ProbeActions = "TERM"
	EnvTermProgram ProbeActions = "TERM_PROGRAM"
	OS             ProbeActions = "OS"
)

type TerminalMD struct {
	ID       string
	Kind     TerminalKind      // According to Terminfo
	reports  []TermInfoReports // According to Terminfo
	strategy []ProbeActions    // According to notcurses & testing
	probe    func(ProbeData) (bool, error)
}

type ProbeData struct {
	OS             string
	XtermVersion   string
	XtermGetTcap   string
	DA1Pp          string
	DA1Ps          string
	DA2Pp          string
	DA2Pv          string
	DA2Pc          string
	DA3            string
	EnvTerm        string
	EnvTermProgram string
}
