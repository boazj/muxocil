package common

type Session struct { // Like Session of teamocil
	Name    string    `yaml:"name"` // You won't have a session name in all cases
	Windows []*Window `yaml:"windows"`
}

type Window struct {
	Name   string    `yaml:"name"`
	Root   string    `yaml:"root"`
	Layout MuxLayout `yaml:"layout"`
	Panes  []*Pane   `yaml:"panes"`
	Focus  bool      `yaml:"focus"`

	Command  string   `yaml:"command"`  // Backwards compatibility with itermocil //TODO: consider if needed
	Commands []string `yaml:"commands"` // Backwards compatibility with itermocil //TODO: consider if needed

	Options map[string]string `yaml:"options"`
}

type Pane struct {
	Commands []string `yaml:"commands"`
	Focus    bool     `yaml:"focus"`
}

type MuxLayout string

const (
	Tiled                  MuxLayout = "tiled"
	EvenHorizontal         MuxLayout = "even-horizontal"
	EvenVertical           MuxLayout = "even-vertical"
	MainVertical           MuxLayout = "main-vertical"
	MainVerticalFlipped    MuxLayout = "main-vertical-flipped"
	MainVerticalMirrored   MuxLayout = "main-vertical-mirrored"
	MainHorizontal         MuxLayout = "main-horizontal"
	MainHorizontalFlipped  MuxLayout = "main-horizontal-flipped"
	MainHorizontalMirrored MuxLayout = "main-horizontal-mirrored"
	ThreeColumns           MuxLayout = "3_columns"
	DoubleMainHorizontal   MuxLayout = "double-main-horizontal"
	DoubleMainVertical     MuxLayout = "double-main-vertical"
)
