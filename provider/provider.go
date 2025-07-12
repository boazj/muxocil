// Package provider is the infrastructure for all multiplers and emulator providers
package provider

import (
	"fmt"
	"runtime"
	"slices"

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
}

var (
	Providers map[common.MuxID]Mux
	hints     struct {
		muxEnv      map[string]common.MuxID
		terminalEnv map[string]common.MuxID
		programEnv  map[string]common.MuxID
	}
)

func init() {
	Providers = map[common.MuxID]Mux{
		common.Tmux:    {common.Tmux, "tmux", common.Multiplexer, []common.OS{common.Windows, common.MacOS, common.Linux}},
		common.Zellij:  {common.Zellij, "Zellij", common.Multiplexer, []common.OS{common.Windows, common.MacOS, common.Linux}},
		common.Iterm2:  {common.Iterm2, "iTerm2", common.Emulator, []common.OS{common.MacOS}},
		common.Kitty:   {common.Kitty, "Kitty", common.Emulator, []common.OS{common.MacOS, common.Linux}},
		common.Wezterm: {common.Wezterm, "WezTerm", common.Emulator, []common.OS{common.Windows, common.MacOS, common.Linux}},
	}

	hints.muxEnv = map[string]common.MuxID{
		"TMUX":                common.Tmux,
		"TMUX_PANE":           common.Tmux,
		"ZELLIJ":              common.Zellij,
		"ZELLIJ_SESSION_NAME": common.Zellij,
		"ITERM_SESSION_ID":    common.Iterm2,
		"KITTY_WINDOW_ID":     common.Kitty,
		"WEZTERM_EXECUTABLE":  common.Wezterm,
	}

	hints.terminalEnv = map[string]common.MuxID{
		"xterm-kitty": common.Kitty,
	}

	hints.programEnv = map[string]common.MuxID{
		"iTerm.app": common.Iterm2,
		"WezTerm":   common.Wezterm,
	}
}

func NewProvider(mux Mux, opts *common.CommandOpts) (Provider, error) {
	switch mux.ID {
	case common.Iterm2:
		p, err := iterm2.NewIterm2(opts)
		return p, utils.Wrap(err, "failed to instantiate provider")
	case common.Kitty:
		p, err := kitty.NewKitty(opts)
		return p, utils.Wrap(err, "failed to instantiate provider")
	case common.Tmux:
		p, err := tmux.NewTmux(opts)
		return p, utils.Wrap(err, "failed to instantiate provider")
	case common.Wezterm:
		p, err := wezterm.NewWezterm(opts)
		return p, utils.Wrap(err, "failed to instantiate provider")
	case common.Zellij:
		p, err := zellij.NewZellij(opts)
		return p, utils.Wrap(err, "failed to instantiate provider")
	}
	return nil, fmt.Errorf("unsupported multiplexer identifier: %s", mux.ID)
}

func FromEnv(opts *common.CommandOpts) (Provider, error) {
	exclude := make([]common.MuxID, 0)
	if runtime.GOOS != "darwin" {
		exclude = append(exclude, common.Iterm2)
	}
	if runtime.GOOS == "windows" {
		exclude = append(exclude, common.Kitty)
	}
	term := utils.GetEnvOr("OVERRIDE_TERM", "TERM")
	program := utils.GetEnvOr("OVERRIDE_TERM_PROGRAM", "TERM_PROGRAM")

	emu, ok := hints.terminalEnv[term]
	if ok && !slices.Contains(exclude, emu) {
		p, err := NewProvider(Providers[emu], opts)
		return p, err
	}
	emu, ok = hints.programEnv[program]
	if ok && !slices.Contains(exclude, emu) {
		p, err := NewProvider(Providers[emu], opts)
		return p, err
	}

	for k, v := range hints.muxEnv {
		if !slices.Contains(exclude, emu) && utils.IsEnvExists(k) {
			p, err := NewProvider(Providers[v], opts)
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

func Process(p Provider, layoutPath string) error {
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

type Provider interface {
	GetID() common.MuxID
	CreateSession(session *common.Session) error
	CreateWindow(window *common.Window, index int) error
	CreatePane(window *common.Window, pane *common.Pane, index int) error

	GetCommads() []string
}
