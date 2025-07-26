package iterm2

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/boazj/muxocil/common"
	"github.com/boazj/muxocil/utils"
	"github.com/charmbracelet/log"
)

// Get version of iTerm. 'iTerm2' (iTerm 2.9+) has better API.
func (t *Iterm2) getVersion() (int, int, int, error) {
	// TODO: deal with beta and nightly
	cmd := exec.Command(
		"osascript",
		"-e",
		fmt.Sprintf("get version of application \"%s\"", AppName),
	)
	out, err := cmd.Output()
	if err != nil {
		return -1, -1, -1, fmt.Errorf("could not fetch iTerm2 version: %v", err)
	}

	version := strings.Split(strings.TrimSpace(string(out)), ".")
	if len(version) != 3 {
		return -1, -1, -1, fmt.Errorf("could not parse iTerm2 version: %s", version)
	}
	major, err1 := strconv.Atoi(version[0])
	minor, err2 := strconv.Atoi(version[1])
	micro, err3 := strconv.Atoi(version[2])
	if err = errors.Join(err1, err2, err3); err != nil {
		return -1, -1, -1, fmt.Errorf("could not parse iTerm2 version: %v", err)
	}
	return major, minor, micro, nil
}

// Get the number of panes already existing in the current window. This is used only for old iTerm.
func (t *Iterm2) getNumPanesInCurrentWindow() (int, error) {
	cmd := exec.Command(
		"osascript",
		"-e",
		fmt.Sprintf("tell application \"%s\" to count sessions of current terminal", AppName),
	)
	out, err := cmd.Output()
	if err != nil {
		return -1, fmt.Errorf("could not fetch number of panes in current window: %v", err)
	}

	value, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return -1, fmt.Errorf("could not parse number of panes in current window: %v", err)
	}
	return value, nil
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

// Once we have laid out the panes we need, we can now navigate
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
		tellTarget = fmt.Sprintf("%s session of current terminal", utils.Ordinal(pane))
	}

	// Setting the pane name is mercifully the same across both iTerm versions.
	if name != "" {
		nameCommand = fmt.Sprintf("set name to \"%s\"", name)
	}
	command := strings.Join(commands, "; ")
	return SingletonScript(Aprintf(TellWrite, tellTarget, command, nameCommand))
}

// Runs the list of commands in the current pane
//
//lint:ignore U1000 in dev
func (t *Iterm2) initiateWindow(commands []string) *AppleScript {
	command := strings.Join(commands, "; ")
	return SingletonScript(Aprintf(TellWrite, TargetSessionOfCurWindow, command, ""))
}

// Switch focus to the specified pane
//
//lint:ignore U1000 in dev
func (t *Iterm2) focusOnPane(paneIndex int) *AppleScript {
	if paneIndex == -1 {
		return &AppleScript{}
	}
	if !t.newIterm && !t.here {
		paneIndex--
	}

	// Determine the correct target for Applescript's 'tell' command based upon iTerm version.
	if t.newIterm {
		return SingletonScript(Aprintf(TellPaneSelected, paneIndex))
	} else {
		return newAppleScript().Append(utils.Times(paneIndex-1, selectNextPane())...)
	}
}
