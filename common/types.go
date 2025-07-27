package common

type CommandOpts struct {
	Layout Session

	Here  bool
	Show  bool
	Debug bool
}

// TODO: force reuse existing layout (session)
// TODO: understand if there is a real use-case for precommand (running command before any other thing, inside of terminal) (itermocil BC)

type MuxID string

//go:generate stringer -type ProviderType
type ProviderType int

type OS string

const (
	Iterm2  MuxID = "iterm2"
	Kitty   MuxID = "kitty"
	Tmux    MuxID = "tmux"
	Wezterm MuxID = "wezterm"
	Zellij  MuxID = "zellij"
)

const (
	Multiplexer ProviderType = iota
	Emulator
)

const (
	Windows OS = "Windows"
	MacOS   OS = "macOS"
	Linux   OS = "Linux"
)

type Provider interface {
	GetID() MuxID
	CreateSession(session *Session) error
	CreateWindow(window *Window, index int) error
	CreatePane(window *Window, pane *Pane, index int) error

	GetCommads() []string
}
