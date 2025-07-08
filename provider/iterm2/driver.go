package iterm2

import (
	"fmt"
	"os"
	"strings"

	"github.com/boazj/muxocil/common"
	"github.com/boazj/muxocil/utils"
	"github.com/charmbracelet/log"
)

const asInitiatePane = ` tell %s
  write text \"%s\"
  %s
 end tell
`

const asInitiateWindow = ` tell current session of current window
  write text "%s"
 end tell
`

const asFocusOnPaneNewIterm = ` tell pane_%d
  select
 end tell
`

func (t *Iterm2) getNumPanesInCurrentWindow() {
	// TODO:
}

// Create a set of Applescript instructions to generate the desired
// layout of panes. Attempt to match teamocil layout behaviour as
// closely as is possible.
//
// See 'arrangePanesOldIterm' for an alternate version for
// generating a version for old iTerm.
func (t *Iterm2) arrangePanes(panes int, layout common.MuxLayout) *AppleScript {
	win := newNewItermWindow(panes, layout)
	ac, err := win.arrangePanes()
	if err != nil {
		log.Fatal("unknown layout strategy ", "layout", layout)
		os.Exit(common.ExitProviderUnknownLayout)
	}
	return ac
}

// Create a set of Applescript instructions to generate the desired
// layout of panes. Attempt to match teamocil layout behaviour as
// closely as is possible.
//
// See 'arrangePanes' for the main version used for generating
// the script for the newer iTerm.
func (t *Iterm2) arrangePanesOldIterm(panes int, layout common.MuxLayout) *AppleScript {
	win := newOldItermWindow(panes, layout)
	ac, err := win.arrangePanes()
	if err != nil {
		log.Fatal("unknown layout strategy ", "layout", layout)
		os.Exit(common.ExitProviderUnknownLayout)
	}
	return ac
}

// Once we have layed out the panes we need, we can now navigate
// to the specified starting directory and run the specified
// commands for each pane.
//
//lint:ignore U1000 in dev
func (t *Iterm2) initiatePane(pane int, commands []string, name string) *AppleScript {
	var tellTarget string
	var nameCommand string

	// Determine the correct target for Applescript's 'tell' command based upon iTerm version.
	if t.newIterm {
		tellTarget = fmt.Sprintf("pane_%d", pane)
	} else {
		// Converts numbers to 2nd, 3rd, 4th format for Applescript
		tellTarget = fmt.Sprintf("%s session of current terminal", utils.Ordinal(pane))
	}

	// Setting the pane name is mercifully the same across both iTerm versions.
	if name != "" {
		nameCommand = fmt.Sprintf("set name to \"%s\"", name)
	}
	command := strings.Join(commands, "; ")
	return SingletonScript(AsCmd(fmt.Sprintf(asInitiatePane, tellTarget, command, nameCommand)))
}

// Runs the list of commands in the current pane
//
//lint:ignore U1000 in dev
func (t *Iterm2) initiateWindow(commands []string) *AppleScript {
	command := strings.Join(commands, "; ")
	return SingletonScript(AsCmd(fmt.Sprintf(asInitiateWindow, command)))
}

// Switch focus to the specified pane
//
//lint:ignore U1000 in dev
func (t *Iterm2) focusOnPane(paneIndex int) *AppleScript {
	if paneIndex == -1 {
		return &AppleScript{}
	}
	if !t.newIterm && !t.here {
		paneIndex = paneIndex - 1
	}

	// Determine the correct target for Applescript's 'tell' command based upon iTerm version.
	if t.newIterm {
		return SingletonScript(AsCmd(fmt.Sprintf(asFocusOnPaneNewIterm, paneIndex)))
	} else {
		return newAppleScript().Append(utils.Times(paneIndex-1, selectNextPane())...)
	}
}
