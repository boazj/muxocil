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

type (
	MuxId     string
	Processor func(opts *common.CommandOpts) error
)

const (
	Iterm2  MuxId = "iterm2"
	Kitty   MuxId = "kitty"
	Tmux    MuxId = "tmux"
	Wezterm MuxId = "wezterm"
	Zellij  MuxId = "zellij"
)

var muxEnvHints = map[string]MuxId{
	"TMUX":                Tmux,
	"TMUX_PANE":           Tmux,
	"ZELLIJ":              Zellij,
	"ZELLIJ_SESSION_NAME": Zellij,
	"ITERM_SESSION_ID":    Iterm2,
	"KITTY_WINDOW_ID":     Kitty,
	"WEZTERM_EXECUTABLE":  Wezterm,
}

var terminalEnvTermHints = map[string]MuxId{
	"xterm-kitty": Kitty,
}

var terminalEnvProgramHints = map[string]MuxId{
	"iTerm.app": Iterm2,
	"WezTerm":   Wezterm,
}

func NewProvider(id MuxId, opts *common.CommandOpts) (Provider, error) {
	switch id {
	case Iterm2:
		p, err := iterm2.NewIterm2(opts)
		return p, err
	case Kitty:
		p, err := kitty.NewKitty(opts)
		return p, err
	case Tmux:
		p, err := tmux.NewTmux(opts)
		return p, err
	case Wezterm:
		p, err := wezterm.NewWezterm(opts)
		return p, err
	case Zellij:
		p, err := zellij.NewZellij(opts)
		return p, err
	}
	return nil, fmt.Errorf("unsupported multiplexer identifier: %s", id)
}

func FromEnv(opts *common.CommandOpts) (Provider, error) {
	exclude := make([]MuxId, 0)
	if runtime.GOOS != "darwin" {
		exclude = append(exclude, Iterm2)
	}
	if runtime.GOOS == "windows" {
		exclude = append(exclude, Kitty)
	}
	term := utils.GetEnvOr("OVERRIDE_TERM", "TERM")
	program := utils.GetEnvOr("OVERRIDE_TERM_PROGRAM", "TERM_PROGRAM")

	emu, ok := terminalEnvTermHints[term]
	if ok && !slices.Contains(exclude, emu) {
		p, err := NewProvider(emu, opts)
		return p, err
	}
	emu, ok = terminalEnvProgramHints[program]
	if ok && !slices.Contains(exclude, emu) {
		p, err := NewProvider(emu, opts)
		return p, err
	}

	for k, v := range muxEnvHints {
		if !slices.Contains(exclude, emu) && utils.IsEnvExists(k) {
			p, err := NewProvider(v, opts)
			return p, err
		}
	}

	return nil, fmt.Errorf("cant recognize terminal emulator or multiplexer via env")
}

func NormalizeValidate(layout *common.Layout) error {
	if len(layout.Windows) == 0 {
		return fmt.Errorf("layout must have at least one window defined")
	}

	focusedWindows := gofn.Filter(layout.Windows, func(win common.Window) bool {
		return win.Focus
	})
	if len(focusedWindows) > 1 {
		return fmt.Errorf("multiple focused windows found")
	} else if len(focusedWindows) == 0 {
		layout.Windows[len(layout.Windows)-1].Focus = true // If no window is marked for focus - mark the last one
	}

	if len(gofn.Filter(layout.Windows, func(win common.Window) bool {
		return len(gofn.Filter(win.Panes, func(pane common.Pane) bool {
			return pane.Focus
		})) > 1
	})) > 1 {
		return fmt.Errorf("multiple focused panes found in the same window")
	}

	for _, win := range layout.Windows {
		if len(win.Panes) > 0 {
			// panes before command and commands
			win.Command = ""
			win.Commands = []string{}
		} else if win.Command != "" {
			// command before commands, normalizing into the commands array to simplify code later
			win.Commands = []string{win.Command}
			win.Command = "" // to avoid confusion
		}

		if len(gofn.Filter(win.Panes, func(pane common.Pane) bool {
			return pane.Focus
		})) == 0 {
			win.Panes[len(win.Panes)-1].Focus = true // If no pane within the window is marked for focus - mark the last one
		}
	}

	return nil
}

func Proccessor(opts *common.CommandOpts) error {
	p, err := FromEnv(opts)
	if err != nil {
		return fmt.Errorf("cannot load provider for current multiplexer: %v", err)
	}

	// TODO: load yml
	var layout *common.Layout = nil
	err = NormalizeValidate(layout)
	if err != nil {
		return fmt.Errorf("encountered issue validating yaml layout: %v", err)
	}
	if layout.Name != "" {
		p.CreateLayout(layout)
	}
	for i, win := range layout.Windows {
		p.CreateWindow(&win, i)
		for j, pane := range win.Panes {
			p.CreatePane(&win, &pane, j)
		}
	}
	return nil
}

type Provider interface {
	CreateLayout(layout *common.Layout) error
	CreateWindow(window *common.Window, index int) error
	CreatePane(window *common.Window, pane *common.Pane, index int) error

	GetCommads() []string
}
