package common

type Layout struct { // Like Session of teamocil
	Name    string   `yaml:"name"` // You won't have a session name in all cases
	Windows []Window `yaml:"windows"`
}

type Window struct {
	Name   string `yaml:"name"`
	Root   string `yaml:"root"`
	Layout string `yaml:"layout"`
	Panes  []Pane `yaml:"panes"`
	Focus  bool   `yaml:"focus"`

	Command  string   `yaml:"command"`  // Backwards compatibility with itermocil //TODO: consider if needed
	Commands []string `yaml:"commands"` // Backwards compatibility with itermocil //TODO: consider if needed

	Options map[string]string `yaml:"options"`
}

type Pane struct {
	Commands []string `yaml:"commands"`
	Focus    bool     `yaml:"focus"`
}
