package termprobe

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
	probe    func(ProbeData) (bool, error)
}

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
