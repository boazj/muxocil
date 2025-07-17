// Package zellij represents the provider for the Zellij Multiplexr
package zellij

import (
	"github.com/boazj/muxocil/common"
	"github.com/boazj/muxocil/utils"
)

type Zellij struct{}

func NewProvider(opts *common.CommandOpts) (common.Provider, error) {
	return &Zellij{}, nil
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

func Detect(cfg *common.Config) (bool, bool) {
	inZellij := cfg.Hints.Zellij != "" && cfg.Hints.ZellijSessionName != ""
	hasZellij := utils.CheckIfCmdInPath("zellij")
	return inZellij, hasZellij
}
