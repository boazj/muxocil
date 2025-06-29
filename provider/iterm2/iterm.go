package iterm2

import "github.com/boazj/muxocil/common"

type Iterm2 struct{}

func NewIterm2(opts *common.CommandOpts) (*Iterm2, error) {
	return nil, nil
}

func (t *Iterm2) CreateLayout(layout *common.Layout) error {
	// TODO: impl
	return nil
}

func (t *Iterm2) CreateWindow(window *common.Window, index int) error {
	// TODO: impl
	return nil
}

func (t *Iterm2) CreatePane(window *common.Window, pane *common.Pane, index int) error {
	// TODO: impl
	return nil
}

func (t *Iterm2) GetCommads() []string {
	// TODO: impl
	return nil
}
