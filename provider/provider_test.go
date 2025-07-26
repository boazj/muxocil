package provider

import (
	"strings"
	"testing"

	"github.com/boazj/muxocil/common"
)

// cfg := &common.Config{
// 	Hints: common.EnvHints{
// 		Tmux:              "/tmp/tmux-1000/default,27349,0",
// 		TmuxPane:          "%16",
// 		Zellij:            "",
// 		ZellijSessionName: "",
// 		ItermSessionID:    "",
// 		KittyWindowID:     "",
// 		WeztermExecutable: "",
// 	},
//
// 	Term:                "xterm-256-color",
// 	OverrideTerm:        "",
// 	TermProgram:         "WezTerm",
// 	OverrideTermProgram: "",
//
// 	LayoutSearchLocations:          []string{"$HOME/.teamocil"},
// 	LayoutSearchExcludeLocations:   []string{},
// 	Editor:                         "vim",
// 	OverrideEditor:                 "",
// 	UseProvider:                    "",
// 	UseProviderFailFast:            true,
// 	ValidateProviderFromEnv:        true,
// 	InfferProviderFromEnv:          true,
// 	LaunchMultiplexerApp:           false,
// 	ShowPaneTitle:                  true,
// 	PaneTitleReplaceWhitespace:     false,
// 	PaneTitleReplaceWhitespaceChar: "_",
// 	IgnoreWindowLayout:             false,
// 	ReplaceIfSessionExists:         true,
// }

func hTmux() *common.EnvHints {
	return &common.EnvHints{Tmux: "/tmp/tmux-1000/default,27349,0", TmuxPane: "%1"}
}

func hZellij() *common.EnvHints {
	return &common.EnvHints{Zellij: "0", ZellijSessionName: "name"}
}

func hWezterm() *common.EnvHints {
	return &common.EnvHints{}
}

func hIterm2() *common.EnvHints {
	return &common.EnvHints{ItermSessionID: "w0t0p0"}
}

func hKitty() *common.EnvHints {
	return &common.EnvHints{KittyWindowID: "id"}
}

func conf(term string, prog string, oprog string, launch bool, hints *common.EnvHints) *common.Config {
	return &common.Config{
		Hints:                *hints,
		Term:                 term,
		OverrideTerm:         "",
		TermProgram:          prog,
		OverrideTermProgram:  oprog,
		LaunchMultiplexerApp: launch,
	}
}

const (
	XTERM  = "xterm-256-color"
	SCREEN = "screen"
	TCOLOR = "tmux-256color"

	TMUX    = "tmux"
	WEZTERM = "WezTerm"
	ITERM   = "iTerm.app"
	KITTY   = "xterm-kitty"
)

func TestFromEnv(t *testing.T) {
	opts := &common.CommandOpts{}
	tests := []struct {
		name string
		in   *common.Config
		w    common.MuxID
		werr string
	}{
		{"TmuxXterm", conf(XTERM, TMUX, "", false, hTmux()), common.Tmux, ""},
		{"TmuxScreen", conf(SCREEN, TMUX, "", false, hTmux()), common.Tmux, ""},
		{"TmuxTColor", conf(TCOLOR, TMUX, "", false, hTmux()), common.Tmux, ""},
		{"TmuxWithoutHints", conf(
			TCOLOR,
			TMUX,
			"",
			false,
			&common.EnvHints{},
		), "", "cant recognize terminal emulator or multiplexer via env"},
		{"TmuxTermZellijHints", conf(TCOLOR, TMUX, "", false, hZellij()), common.Zellij, ""},
		{"Zellij", conf(XTERM, WEZTERM, "", false, hZellij()), common.Zellij, ""},
		{"iTerm2", conf(XTERM, ITERM, "", false, hIterm2()), common.Iterm2, ""},
		{"WezTerm", conf(XTERM, WEZTERM, "", false, hWezterm()), common.Wezterm, ""},
		{"Kitty", conf(KITTY, "", "", false, hKitty()), common.Kitty, ""},
		// FIXME: term_program is not set properly for kitty by design,
		// so this is recognized as wezterm, nit to use ansi CSI escape code to better recognize apps
		// {"Kitty", conf(KITTY, "", WEZTERM, "", false, hKitty()), common.Kitty, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans, err := FromEnv(tt.in, opts)
			if tt.werr == "" && err == nil && ans.GetID() != tt.w {
				t.Errorf("got %s, want %s", ans.GetID(), tt.w)
			}
			if tt.werr != "" && !strings.Contains(err.Error(), tt.werr) {
				t.Errorf("got %v, want %s", err, tt.werr)
			}
		})
	}
}
