package iterm2

import (
	"fmt"
	"math"
	"strings"

	"github.com/boazj/muxocil/utils"
)

type ItermLayout string

const (
	Tiled                ItermLayout = "tiled"
	EvenHorizontal       ItermLayout = "even-horizontal"
	EvenVertical         ItermLayout = "even-vertical"
	MainVertical         ItermLayout = "main-vertical"
	MainVerticalFlipped  ItermLayout = "main-vertical-flipped"
	MainHorizontal       ItermLayout = "main-horizontal"
	ThreeColumns         ItermLayout = "3_columns"
	DoubleMainHorizontal ItermLayout = "double-main-horizontal"
	DoubleMainVertical   ItermLayout = "double-main-vertical"
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

func splitVertically() string {
	return pressKeystroke("d")
}

func splitHorizontally() string {
	return pressKeystroke("D")
}

func selectPreviousPane() string {
	return pressKeystroke("]")
}

func selectNextPane() string {
	return pressKeystroke("[")
}

func pressKeystroke(key string) string {
	// TODO: modifiers
	return fmt.Sprintf("tell i term application \"System Events\" to keystroke \"%s\" using command down", key)
}

func pressKeyCode(key string) string {
	// TODO: modifiers
	return fmt.Sprintf("tell i term application \"System Events\" to key code \"%s\" using {command down, option down}", key)
}

// Create a set of Applescript instructions to generate the desired
// layout of panes. Attempt to match teamocil layout behaviour as
// closely as is possible.
//
// See 'arrange_panes' for the main version used for generating
// the script for the newer iTerm.
//
//lint:ignore U1000 in dev
func (t *Iterm2) arrangePanesOldIterm(numPanes int, layout ItermLayout) ([]string, error) {
	// TODO: this entire crappy function relies on having the defualt hotkeys
	if layout == "" {
		layout = Tiled
	}
	as := make([]string, 0)

	// If we have just one pane we don't need to do any splitting.
	if numPanes <= 1 {
		return []string{}, nil
	}

	// tmux seems to treat the first 2 tiles of a tiled layout like this
	if numPanes == 2 && layout == Tiled {
		layout = EvenVertical
	}
	if layout == EvenHorizontal {
		// 'even-horizontal' layouts just split vertically across the screen
		as = utils.AppendTimes(as, numPanes-1, splitVertically())
		as = append(as, selectPreviousPane()) // Focus back on the first pane

	} else if layout == EvenVertical {
		// 'even-vertical' layouts just split horizontally down the screen
		as = utils.AppendTimes(as, numPanes-1, splitVertically())
		as = append(as, selectPreviousPane()) // Focus back on the first pane

	} else if layout == MainVertical {
		// 'main-vertical' layouts have one left pane that is full height,
		// and then split the remaining panes horizontally down the right
		as = append(as, splitVertically())
		as = utils.AppendTimes(as, numPanes-2, splitHorizontally())
		as = append(as, selectPreviousPane()) // Focus back on the first pane

	} else if layout == MainVerticalFlipped {
		// 'main-vertical-flipped' layouts have one right pane that is full height,
		// and then split the remaining panes horizontally down the left
		as = append(as, splitVertically())
		as = append(as, pressKeystroke("[")) // Focus back on the first pane
		as = utils.AppendTimes(as, numPanes-2, splitHorizontally())

	} else if layout == MainHorizontal {
		// 'main-horizontal' layouts have one left pane  that is full height,
		// and then split the remaining panes horizontally down the right
		as = append(as, splitHorizontally())
		as = utils.AppendTimes(as, numPanes-2, splitVertically())
		as = append(as, selectPreviousPane()) // Focus back on the first pane

	} else if layout == Tiled {
		// 'tiled' layouts create 2 columns and then however many rows as
		// needed. If there are odd number of panes then the bottom pane
		// spans two columns. Panes are numbered top to bottom, left to right.
		verticalSplits := int(math.Ceil(float64(numPanes)/2)) - 1
		secondColumns := numPanes / 2

		as = utils.AppendTimes(as, verticalSplits, splitHorizontally())

		if verticalSplits > 0 {
			// If we split vertically at all then move 'down' a pane to take
			// us back to the initial pane.
			as = append(as, pressKeyCode("125"))
		}

		as = utils.AppendTimes(
			as,
			secondColumns,
			splitVertically(),
			selectPreviousPane(),
		)

		if numPanes%2 != 0 {
			as = append(as, selectPreviousPane())
		}
	} else if layout == ThreeColumns {
		// '3_columns' layouts create 3 columns and then however many rows as
		// needed. If there are odd number of panes then the bottom pane
		// spans two columns. Panes are numbered top to bottom, left to right.
		verticalSplits := int(math.Ceil(float64(numPanes)/3)) - 1
		i := 1 + verticalSplits
		as = utils.AppendTimes(as, verticalSplits, splitHorizontally())

		re := numPanes - i
		if re%2 == 0 {
			as = utils.AppendTimes(as, re/2, selectPreviousPane(), splitVertically(), splitVertically())
		} else {
			as = utils.AppendTimes(as, (re-1)/2, selectPreviousPane(), splitVertically(), splitVertically())
			as = append(as, selectPreviousPane())
			as = append(as, splitVertically())
		}

		// as = append(as, selectPreviousPane())
		// i = i + 1
		// as = append(as, splitVertically())
		// if i >= numPanes {
		// 	break
		// }
		// i = i + 1
		// as = append(as, splitVertically())
		// if i >= numPanes {
		// 	break
		// }

		//
		// while True:
		//     self.applescript.append(prefix + 'keystroke "]" using command down')
		//     i += 1
		//     self.applescript.append(prefix + 'keystroke "d" using command down')
		//     if i >= num_panes:
		//         break
		//     i += 1
		//     self.applescript.append(prefix + 'keystroke "d" using command down')
		//     if i >= num_panes:
		//         break
	} else if layout == DoubleMainHorizontal {
		// 'double-main-horizontal' layouts have two left panes that are full height,
		// and then split the remaining panes horizontally down the right

		as = append(as, splitVertically())
		if numPanes > 2 {
			as = append(as, splitVertically())
		}
		if numPanes > 3 {
			as = utils.AppendTimes(as, numPanes-3, splitHorizontally())
		}
	} else if layout == DoubleMainVertical {
		// 'double-main-vertical' layouts have two bottom panes that split the width
		// and then split the remaining panes vertically across the top

		as = append(as, splitHorizontally())
		as = append(as, splitVertically())
		as = append(as, selectPreviousPane())
		if numPanes > 3 {
			as = utils.AppendTimes(as, numPanes-3, splitVertically())
		}
	} else {
		// Raise an exception if we don't recognise the layout setting.
		return []string{}, fmt.Errorf("unknown layout setting")
	}

	// This is all keystroke based and thus takes a moment to happen,
	// so unfortunately (for old iTerm) we have to wait a moment to
	// give all that time to happen.
	as = append(as, "delay 2")
	return as, nil
}

// Once we have layed out the panes we need, we can now navigate
// to the specified starting directory and run the specified
// commands for each pane.
//
//lint:ignore U1000 in dev
func (t *Iterm2) initiatePane(pane int, commands []string, name string) string {
	var tellTarget string
	var nameCommand string

	// Determine the correct target for Applescript's 'tell' command
	// based upon iTerm version.
	if t.newIterm {
		tellTarget = fmt.Sprintf("pane_%d", pane)
	} else {
		// Converts numbers to 2nd, 3rd, 4th format for Applescript
		tellTarget = fmt.Sprintf("%s session of current terminal", utils.Ordinal(pane))
	}

	// Setting the pane name is mercifully the same across both
	// iTerm versions.
	if name != "" {
		nameCommand = fmt.Sprintf("set name to \"%s\"", name)
	}
	// Turn commands list into a string command
	command := strings.Join(commands, "; ")
	// Build the applescript snippet.
	return fmt.Sprintf(asInitiatePane, tellTarget, command, nameCommand)
}

// Runs the list of commands in the current pane
//
//lint:ignore U1000 in dev
func (t *Iterm2) initiateWindow(commands []string) string {
	command := strings.Join(commands, "; ")
	return fmt.Sprintf(asInitiateWindow, command)
}

// Switch focus to the specified pane
//
//lint:ignore U1000 in dev
func (t *Iterm2) focusOnPane(pane int) []string {
	if pane == -1 {
		return []string{}
	}
	as := make([]string, 0)

	if !t.newIterm && !t.here {
		pane = pane - 1
	}

	// Determine the correct target for Applescript's 'tell' command
	// based upon iTerm version.
	if t.newIterm {
		as = append(as, fmt.Sprintf(asFocusOnPaneNewIterm, pane))
	} else {
		as = utils.AppendTimes(as, pane-1, selectPreviousPane())
	}
	return as
}
