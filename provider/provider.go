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

type (
	MuxID     string
	Processor func(opts *common.CommandOpts) error
)

const (
	Iterm2  MuxID = "iterm2"
	Kitty   MuxID = "kitty"
	Tmux    MuxID = "tmux"
	Wezterm MuxID = "wezterm"
	Zellij  MuxID = "zellij"
)

type ProviderType string

const (
	Multiplexer = "Multiplexer"
	Emulator    = "Emulator"
)

type OS string

const (
	Windows = "Windows"
	MacOS   = "macOS"
	Linux   = "Linux"
)

type MuxMD struct {
	ID          MuxID
	Display     string
	Kind        ProviderType
	SupportedOs []OS
}

var Providers = []MuxMD{
	{Tmux, "tmux", Multiplexer, []OS{Windows, MacOS, Linux}},
	{Zellij, "Zellij", Multiplexer, []OS{Windows, MacOS, Linux}},
	{Iterm2, "iTerm2", Emulator, []OS{MacOS}},
	{Kitty, "Kitty", Emulator, []OS{MacOS, Linux}},
	{Wezterm, "WezTerm", Emulator, []OS{Windows, MacOS, Linux}},
}

var muxEnvHints = map[string]MuxID{
	"TMUX":                Tmux,
	"TMUX_PANE":           Tmux,
	"ZELLIJ":              Zellij,
	"ZELLIJ_SESSION_NAME": Zellij,
	"ITERM_SESSION_ID":    Iterm2,
	"KITTY_WINDOW_ID":     Kitty,
	"WEZTERM_EXECUTABLE":  Wezterm,
}

var terminalEnvTermHints = map[string]MuxID{
	"xterm-kitty": Kitty,
}

var terminalEnvProgramHints = map[string]MuxID{
	"iTerm.app": Iterm2,
	"WezTerm":   Wezterm,
}

func NewProvider(id MuxID, opts *common.CommandOpts) (Provider, error) {
	switch id {
	case Iterm2:
		p, err := iterm2.NewIterm2(opts)
		return p, utils.Wrap(err, "failed to instantiate provider")
	case Kitty:
		p, err := kitty.NewKitty(opts)
		return p, utils.Wrap(err, "failed to instantiate provider")
	case Tmux:
		p, err := tmux.NewTmux(opts)
		return p, utils.Wrap(err, "failed to instantiate provider")
	case Wezterm:
		p, err := wezterm.NewWezterm(opts)
		return p, utils.Wrap(err, "failed to instantiate provider")
	case Zellij:
		p, err := zellij.NewZellij(opts)
		return p, utils.Wrap(err, "failed to instantiate provider")
	}
	return nil, fmt.Errorf("unsupported multiplexer identifier: %s", id)
}

func FromEnv(opts *common.CommandOpts) (Provider, error) {
	exclude := make([]MuxID, 0)
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

func Proccessor(opts *common.CommandOpts) error {
	p, err := FromEnv(opts)
	if err != nil {
		return fmt.Errorf("cannot load provider for current multiplexer: %v", err)
	}

	// TODO: load yml
	var session *common.Session = nil
	err = NormalizeValidate(session)
	if err != nil {
		return fmt.Errorf("encountered issue validating yaml layout: %v", err)
	}
	if session.Name != "" {
		p.CreateLayout(session)
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
	CreateLayout(session *common.Session) error
	CreateWindow(window *common.Window, index int) error
	CreatePane(window *common.Window, pane *common.Pane, index int) error

	GetCommads() []string
}
