// Package wezterm represents the provider for the WezTerm Terminal Emulator
package wezterm

import "github.com/boazj/muxocil/common"

type Wezterm struct{}

func NewWezterm(opts *common.CommandOpts) (*Wezterm, error) {
	return nil, nil
}

func (w *Wezterm) CreateSession(session *common.Session) error {
	// TODO: impl
	return nil
}

func (w *Wezterm) CreateWindow(window *common.Window, index int) error {
	// TODO: impl
	return nil
}

func (w *Wezterm) CreatePane(window *common.Window, pane *common.Pane, index int) error {
	// TODO: impl
	return nil
}

func (w *Wezterm) GetCommads() []string {
	// TODO: impl
	return nil
}
