// Package iterm2 represents the provider for the iTerm2 Terminal Emulator
package iterm2

import (
	"fmt"

	"github.com/boazj/muxocil/common"
)

const AppName = "iTerm2"

type Iterm2 struct {
	newIterm         bool
	here             bool
	initialPaneCount int

	script AppleScript
}

func (t *Iterm2) getScript() string {
	return t.script.Raw()
}

func NewIterm2(opts *common.CommandOpts) (*Iterm2, error) {
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
	TellActivate  = "tell application \"%s\" to activate"
	TellCreateTab = `tell current window
      create tab with default profile
    end tell
  `
)

func (t *Iterm2) CreateSession(session *common.Session) error {
	t.script.Append(AsCmd(fmt.Sprintf(TellActivate, AppName)))

	// TODO: if I decide to introduce pre, this is the origin
	// if 'pre' in self.parsed_config:
	//   self.applescript.append('do shell script "' + self.parsed_config['pre'] + ';"')

	if !t.here {
		if t.newIterm {
			t.script.Append(TellCreateTab)
		} else {
			t.script.Append(
				"delay 0.3",
				"tell i term application \"System Events\" to keystroke \"t\" using command down",
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
			"tell i term application \"System Events\" to keystroke \"t\" using command down",
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
