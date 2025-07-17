// Package kitty represents the provider for the Kitty Terminal Emulator
package kitty

import (
	"github.com/boazj/muxocil/common"
	"github.com/tiendc/gofn"
)

type Kitty struct{}

func NewProvider(opts *common.CommandOpts) (common.Provider, error) {
	return &Kitty{}, nil
}

func (k *Kitty) GetID() common.MuxID {
	return common.Kitty
}

func (k *Kitty) CreateSession(session *common.Session) error {
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

func Detect(cfg *common.Config) (bool, bool) {
	term := gofn.FirstNonEmpty(cfg.OverrideTerm, cfg.Term)
	inKitty := term == "xterm-kitty" && cfg.Hints.KittyWindowID != ""
	return inKitty, false
}
