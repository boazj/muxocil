// Package zellij represents the provider for the Zellij Multiplexr
package zellij

import "github.com/boazj/muxocil/common"

type Zellij struct{}

func NewProvider(opts *common.CommandOpts) (common.Provider, error) {
	return nil, nil
}

func (z *Zellij) GetID() common.MuxID {
	return common.Zellij
}

func (z *Zellij) CreateSession(session *common.Session) error {
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
