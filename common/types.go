package common

type CommandOpts struct {
	Layout Layout

	Here  bool
	Show  bool
	Debug bool
}

// TODO: force reuse existing layout (session)
// TODO: attach if layout(session) exists - -A on new-session
// TODO: force specific provider
// TODO: activate non-emulator multiplexer (tmux, zellij) if no executed inside (candidate for config file!)
// TODO: understand if there is a real use-case for precommand (running command before any other thing, inside of terminal) (itermocil BC)
// TODO: support dedicated title for panes
type Config struct{}
