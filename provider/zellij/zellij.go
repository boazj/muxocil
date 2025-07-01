// Package zellij represents the provider for the Zellij Multiplexr
package zellij

import "github.com/boazj/muxocil/common"

type Zellij struct{}

func NewZellij(opts *common.CommandOpts) (*Zellij, error) {
	return nil, nil
}

func (z *Zellij) CreateLayout(session *common.Session) error {
	// TODO: impl
	return nil
}

func (z *Zellij) CreateWindow(window *common.Window, index int) error {
	// TODO: impl
	return nil
}

func (z *Zellij) CreatePane(window *common.Window, pane *common.Pane, index int) error {
	// TODO: impl
	return nil
}

func (z *Zellij) GetCommads() []string {
	// TODO: impl
	return nil
}
