package kitty

import "github.com/boazj/muxocil/common"

type Kitty struct{}

func NewKitty(opts *common.CommandOpts) (*Kitty, error) {
	return nil, nil
}

func (k *Kitty) CreateLayout(layout *common.Layout) error {
	// TODO: impl
	return nil
}

func (k *Kitty) CreateWindow(window *common.Window, index int) error {
	// TODO: impl
	return nil
}

func (k *Kitty) CreatePane(window *common.Window, pane *common.Pane, index int) error {
	// TODO: impl
	return nil
}

func (k *Kitty) GetCommads() []string {
	// TODO: impl
	return nil
}
