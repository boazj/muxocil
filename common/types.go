package common

type CommandOpts struct {
	Layout Layout

	Here  bool
	Show  bool
	Debug bool
}

// TODO: force reuse existing layout (session)
// TODO: understand if there is a real use-case for precommand (running command before any other thing, inside of terminal) (itermocil BC)
