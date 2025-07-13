package provider

import (
	"testing"

	"github.com/boazj/muxocil/common"
	"github.com/stretchr/testify/assert"
)

// TODO:
func TestFromEnvBasicTmux(t *testing.T) {
	cfg := &common.Config{
		Hints: common.EnvHints{
			Tmux:              "/tmp/tmux-1000/default,27349,0",
			TmuxPane:          "%16",
			Zellij:            "",
			ZellijSessionName: "",
			ItermSessionID:    "",
			KittyWindowID:     "",
			WeztermExecutable: "",
		},

		Term:                "tmux-256color",
		OverrideTerm:        "",
		TermProgram:         "tmux",
		OverrideTermProgram: "",

		LayoutSearchLocations:          []string{"$HOME/.teamocil"},
		LayoutSearchExcludeLocations:   []string{},
		Editor:                         "vim",
		OverrideEditor:                 "",
		UseProvider:                    "",
		UseProviderFailFast:            true,
		ValidateProviderFromEnv:        true,
		InfferProviderFromEnv:          true,
		LaunchMultiplexerApp:           false,
		ShowPaneTitle:                  true,
		PaneTitleReplaceWhitespace:     false,
		PaneTitleReplaceWhitespaceChar: "_",
		IgnoreWindowLayout:             false,
		ReplaceIfSessionExists:         true,
	}

	opts := &common.CommandOpts{}
	p, err := FromEnv(cfg, opts)
	assert.Nil(t, err)
	assert.Equal(t, common.Tmux, p.GetID(), "")
}
