// Package provider is the infrastructure for all multiplers and emulator providers
package provider

import (
	"fmt"

	"github.com/boazj/muxocil/common"
	"github.com/boazj/muxocil/utils"
	"github.com/tiendc/gofn"
)

var ProviderDefs *providerDefs

func init() {
	ProviderDefs = CreateProviderDefs()
}

func NewProvider(cfg *common.Config, mux common.MuxID, opts *common.CommandOpts) (common.Provider, error) {
	m, ok := ProviderDefs.GetProvider(mux)
	if !ok {
		return nil, fmt.Errorf("unsupported multiplexer identifier: %s", mux)
	}
	p, err := m.constructor(opts)
	return p, utils.Wrap(err, "failed to instantiate provider")
}

//nolint:staticcheck
func FromEnv(cfg *common.Config, opts *common.CommandOpts) (common.Provider, error) {
	muxers := ProviderDefs.GetSupportedMultiplexers()
	emus := ProviderDefs.GetSupportedEmulators()
	// Ensure multiplexers will be checked before emulator multiplexers as they are more specific
	providers := gofn.Concat(muxers, emus)
	for _, v := range providers {
		in, has := v.detector(cfg)
		if has && !in && cfg.LaunchMultiplexerApp {
			// TODO: launch application if present, currently all emus return has == false
			// might be wrong
		}
		if in {
			p, err := NewProvider(cfg, v.ID, opts)
			return p, err
		}
	}

	return nil, fmt.Errorf("cant recognize terminal emulator or multiplexer via env")
}

func NormalizeValidate(session *common.Session) error {
	if len(session.Windows) == 0 {
		return fmt.Errorf("layout must have at least one window defined")
	}

	focusedWindows := gofn.Filter(session.Windows, func(win *common.Window) bool {
		return win.Focus
	})
	if len(focusedWindows) > 1 {
		return fmt.Errorf("multiple focused windows found")
	} else if len(focusedWindows) == 0 {
		session.Windows[len(session.Windows)-1].Focus = true // If no window is marked for focus - mark the last one
	}

	if len(gofn.Filter(session.Windows, func(win *common.Window) bool {
		return len(gofn.Filter(win.Panes, func(pane *common.Pane) bool {
			return pane.Focus
		})) > 1
	})) > 1 {
		return fmt.Errorf("multiple focused panes found in the same window")
	}

	for _, win := range session.Windows {
		if len(win.Panes) > 0 {
			// panes before command and commands
			win.Command = ""
			win.Commands = []string{}
		} else if win.Command != "" {
			// command before commands, normalizing into the commands array to simplify code later
			win.Commands = []string{win.Command}
			win.Command = "" // to avoid confusion
		}

		if len(gofn.Filter(win.Panes, func(pane *common.Pane) bool {
			return pane.Focus
		})) == 0 {
			win.Panes[len(win.Panes)-1].Focus = true // If no pane within the window is marked for focus - mark the last one
		}
	}

	return nil
}

func Process(p common.Provider, layoutPath string) error {
	// TODO: load yml
	var session *common.Session = nil
	err := NormalizeValidate(session)
	if err != nil {
		return fmt.Errorf("encountered issue validating yaml layout: %v", err)
	}
	if session.Name != "" {
		err = p.CreateSession(session)
		if err != nil {
			return fmt.Errorf("encountered issue while creation a session: %v", err)
		}
	}
	for i, win := range session.Windows {
		err = p.CreateWindow(win, i)
		if err != nil {
			return fmt.Errorf("encountered issue while creation a window: %v", err)
		}
		for j, pane := range win.Panes {
			err = p.CreatePane(win, pane, j)
			if err != nil {
				return fmt.Errorf("encountered issue while creation a pane: %v", err)
			}
		}
	}
	return nil
}
