// Package wezterm represents the provider for the WezTerm Terminal Emulator
package wezterm

import (
	"github.com/boazj/muxocil/common"
	"github.com/tiendc/gofn"
)

type Wezterm struct{}

func NewProvider(opts *common.CommandOpts) (common.Provider, error) {
	return &Wezterm{}, nil
}

func (w *Wezterm) GetID() common.MuxID {
	return common.Wezterm
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

func Detect(cfg *common.Config) (bool, bool) {
	program := gofn.FirstNonEmpty(cfg.OverrideTermProgram, cfg.TermProgram)
	inWezterm := program == "WezTerm"
	return inWezterm, false
}
