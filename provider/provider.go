// Package provider is the infrastructure for all multiplers and emulator providers
package provider

import (
	"fmt"
	"runtime"

	"github.com/boazj/muxocil/common"
	"github.com/boazj/muxocil/provider/iterm2"
	"github.com/boazj/muxocil/provider/kitty"
	"github.com/boazj/muxocil/provider/tmux"
	"github.com/boazj/muxocil/provider/wezterm"
	"github.com/boazj/muxocil/provider/zellij"
	"github.com/boazj/muxocil/utils"
	"github.com/tiendc/gofn"
)

type Mux struct {
	ID          common.MuxID
	Display     string
	Kind        common.ProviderType
	SupportedOs []common.OS
	constructor func(*common.CommandOpts) (common.Provider, error)
}

var (
	// TODO: encapsulate
	Providers map[common.MuxID]Mux
	hints     struct {
		muxEnv      map[string]common.MuxID
		terminalEnv map[string]common.MuxID
		programEnv  map[string]common.MuxID
	}
)

func init() {
	Providers = make(map[common.MuxID]Mux)
	Providers[common.Tmux] = Mux{
		common.Tmux,
		"tmux",
		common.Multiplexer,
		[]common.OS{common.Windows, common.MacOS, common.Linux},
		tmux.NewProvider,
	}
	Providers[common.Zellij] = Mux{
		common.Zellij,
		"Zellij",
		common.Multiplexer,
		[]common.OS{common.Windows, common.MacOS, common.Linux},
		zellij.NewProvider,
	}

	if runtime.GOOS == "darwin" {
		Providers[common.Iterm2] = Mux{
			common.Iterm2,
			"iTerm2",
			common.Emulator,
			[]common.OS{common.MacOS},
			iterm2.NewProvider,
		}
	}
	if runtime.GOOS != "windows" {
		Providers[common.Kitty] = Mux{
			common.Kitty,
			"Kitty",
			common.Emulator,
			[]common.OS{common.MacOS, common.Linux},
			kitty.NewProvider,
		}
	}
	Providers[common.Wezterm] = Mux{
		common.Wezterm,
		"WezTerm",
		common.Emulator,
		[]common.OS{common.Windows, common.MacOS, common.Linux},
		wezterm.NewProvider,
	}

	hints.muxEnv = map[string]common.MuxID{
		"Tmux":              common.Tmux,
		"TmuxPane":          common.Tmux,
		"Zellij":            common.Zellij,
		"ZellijSessionName": common.Zellij,
		"ItermSessionID":    common.Iterm2,
		"KittyWindowID":     common.Kitty,
		"WeztermExecutable": common.Wezterm,
	}

	hints.terminalEnv = map[string]common.MuxID{
		"xterm-kitty": common.Kitty,
	}

	hints.programEnv = map[string]common.MuxID{
		// "tmux":      common.Tmux,
		"iTerm.app": common.Iterm2,
		"WezTerm":   common.Wezterm,
	}
}

func NewProvider(cfg *common.Config, mux Mux, opts *common.CommandOpts) (common.Provider, error) {
	m, ok := Providers[mux.ID] // fetching again to avoid tempring
	if !ok {
		return nil, fmt.Errorf("unsupported multiplexer identifier: %s", mux.ID)
	}
	p, err := m.constructor(opts)
	return p, utils.Wrap(err, "failed to instantiate provider")
}

func FromEnv(cfg *common.Config, opts *common.CommandOpts) (common.Provider, error) {
	term := gofn.FirstNonEmpty(cfg.OverrideTerm, cfg.Term)
	program := gofn.FirstNonEmpty(cfg.OverrideTermProgram, cfg.TermProgram)

	emu, ok := hints.terminalEnv[term]
	if ok {
		p, err := NewProvider(cfg, Providers[emu], opts)
		return p, err
	}
	emu, ok = hints.programEnv[program]
	if ok {
		p, err := NewProvider(cfg, Providers[emu], opts)
		return p, err
	}

	for k, v := range processHints(cfg) {
		if v != "" {
			p, err := NewProvider(cfg, Providers[hints.muxEnv[k]], opts)
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
		p.CreateSession(session)
	}
	for i, win := range session.Windows {
		p.CreateWindow(win, i)
		for j, pane := range win.Panes {
			p.CreatePane(win, pane, j)
		}
	}
	return nil
}

func processHints(cfg *common.Config) map[string]string {
	h, _ := utils.AsStringMap(cfg.Hints)
	return h
}
