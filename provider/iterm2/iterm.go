// Package iterm2 represents the provider for the iTerm2 Terminal Emulator
package iterm2

import (
	"github.com/boazj/muxocil/common"
	"github.com/tiendc/gofn"
)

type Iterm2 struct {
	newIterm         bool
	here             bool
	initialPaneCount int

	script AppleScript
}

func (t *Iterm2) getScript() string {
	return t.script.Raw()
}

func NewProvider(opts *common.CommandOpts) (common.Provider, error) {
	t := &Iterm2{
		here:             opts.Here,
		newIterm:         true,
		initialPaneCount: 0,

		script: *newAppleScript(),
	}
	major, minor, _, err := t.getVersion()
	if err != nil {
		return nil, err
	}
	if major < 2 || (major == 2 && minor < 9) {
		t.newIterm = false
	}

	if !t.newIterm && !t.here {
		t.initialPaneCount, err = t.getNumPanesInCurrentWindow()
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

const (
	TellLaunch = `if %s is not running then
      tell %s to launch
    end if
  `
	TellActivate  = "tell application \"%s\" to activate"
	TellCreateTab = `tell current window
      create tab with default profile
    end tell
  `
	TellWrite = `tell %s
      write text \"%s\"
      %s
    end tell
  `
	TellPaneSelected = `tell pane_%d
      select
    end tell
  `

	AppName                  = "iTerm2"
	OldAppName               = "i term"
	TargetApp                = "application iTerm2"
	TargetOldApp             = "i term application"
	TargetCurSession         = "current session"
	TargetCurWindow          = "current window"
	TargetSessionOfCurWindow = TargetCurSession + " of " + TargetCurWindow
)

func (t *Iterm2) GetID() common.MuxID {
	return common.Iterm2
}

func (t *Iterm2) CreateSession(session *common.Session) error {
	t.script.Append(Aprintf(TellActivate, AppName))

	// TODO: if I decide to introduce pre, this is the origin
	// if 'pre' in self.parsed_config:
	//   self.applescript.append('do shell script "' + self.parsed_config['pre'] + ';"')

	if !t.here {
		if t.newIterm {
			t.script.Append(Aprintf(TellLaunch, TargetApp, TargetApp), TellCreateTab)
		} else {
			t.script.Append(
				"delay 0.3",
				pressKeystroke("t", KeyCommand),
				"delay 0.3",
			)
		}
	}

	t.script.Suffix("end tell")
	return nil
}

func (t *Iterm2) CreateWindow(window *common.Window, index int) error {
	// TODO: impl
	if t.newIterm {
		t.script.Append(TellCreateTab)
	} else {
		t.script.Append(
			"delay 0.3",
			pressKeystroke("t", KeyCommand),
			"delay 0.3",
		)
	}
	return nil
}

func (t *Iterm2) CreatePane(window *common.Window, pane *common.Pane, index int) error {
	// TODO: impl
	return nil
}

func (t *Iterm2) GetCommads() []string {
	// TODO: impl
	return nil
}

func Detect(cfg *common.Config) (bool, bool) {
	program := gofn.FirstNonEmpty(cfg.OverrideTermProgram, cfg.TermProgram)
	inIterm2 := program == "iTerm.app" && cfg.Hints.ItermSessionID != ""
	return inIterm2, false
}
