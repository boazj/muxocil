package iterm2

import (
	"fmt"

	"github.com/boazj/muxocil/common"
	"github.com/boazj/muxocil/utils"
)

type oldItermWindow struct {
	panes        int
	createdPanes int

	layout common.MuxLayout

	layoutStrategy map[common.MuxLayout]func()

	script AppleScript
}

func newOldItermWindow(panes int, layout common.MuxLayout) *oldItermWindow {
	if layout == "" {
		layout = common.Tiled
	}

	// tmux seems to treat the first 2 tiles of a tiled layout like this
	if panes == 2 && layout == common.Tiled {
		layout = common.EvenVertical
	}

	win := &oldItermWindow{
		panes:        panes,
		createdPanes: 1, // windows start with an initial pane
		layout:       layout,
		script:       *newAppleScript(),
	}

	strategies := map[common.MuxLayout]func(){
		common.EvenHorizontal:         win.evenHorizontal,
		common.EvenVertical:           win.evenVertical,
		common.MainHorizontal:         win.mainHorizontal,
		common.MainHorizontalMirrored: win.mainHorizontalFlipped,
		common.MainHorizontalFlipped:  win.mainHorizontalFlipped,
		common.MainVertical:           win.mainVertical,
		common.MainVerticalMirrored:   win.mainVerticalFlipped,
		common.MainVerticalFlipped:    win.mainVerticalFlipped,
		common.Tiled:                  win.tiled,
		common.ThreeColumns:           win.threeColumns,
		common.DoubleMainHorizontal:   win.doubleMainHorizontal,
		common.DoubleMainVertical:     win.doubleMainVertical,
	}
	win.layoutStrategy = strategies

	return win
}

// TODO: verify that the mix of splits makes sense, the values seems flipped

func (w *oldItermWindow) arrangePanes() (*AppleScript, error) {
	if w.panes <= 1 {
		// If we have just one pane we don't need to do any splitting.
		return &AppleScript{}, nil
	}
	strategy, ok := w.layoutStrategy[w.layout]
	if !ok {
		return &AppleScript{}, fmt.Errorf("unknown layout setting")
	}
	strategy()

	// This is all keystroke based and thus takes a moment to happen,
	// so unfortunately (for old iTerm) we have to wait a moment to
	// give all that time to happen.
	w.script.Append("delay 2")
	return &w.script, nil
}

// Create a pane to the right of the current pane.
func (w *oldItermWindow) splitVertically() AsCmd {
	w.createdPanes++
	return pressKeystroke("d", KeyCommand)
}

func (w *oldItermWindow) splitVerticallyRemaining() {
	times := w.panes - w.createdPanes
	w.script.Append(utils.Times(times, w.splitVertically())...)
}

// Create a pane below the current pane.
func (w *oldItermWindow) splitHorizontally() AsCmd {
	w.createdPanes++
	return pressKeystroke("D", KeyCommand)
}

func (w *oldItermWindow) splitHorizontallyRemaining() {
	times := w.panes - w.createdPanes
	w.script.Append(utils.Times(times, w.splitHorizontally())...)
}

// 'even-horizontal' layouts just split vertically across the screen.
func (w *oldItermWindow) evenHorizontal() {
	w.splitVerticallyRemaining()
	w.script.Append(selectNextPane())
}

// 'even-vertical' layouts just split horizontally down the screen.
func (w *oldItermWindow) evenVertical() {
	w.splitHorizontallyRemaining()
	w.script.Append(selectNextPane())
}

// 'main-vertical' layouts have one left pane that is full height,
// and then split the remaining panes horizontally down the right.
func (w *oldItermWindow) mainVertical() {
	w.script.Append(w.splitVertically())
	w.splitHorizontallyRemaining()
	w.script.Append(selectNextPane())
}

// 'main-vertical-flipped' layouts have one right pane that is full height,
// and then split the remaining panes horizontally down the left.
func (w *oldItermWindow) mainVerticalFlipped() {
	w.script.Append(w.splitVertically())
	w.script.Append(selectPrevPane())
	w.splitHorizontallyRemaining()
}

// 'main-horizontal' layouts have one upper pane that is full width,
// and then split the remaining panes vertically along the buttom.
func (w *oldItermWindow) mainHorizontal() {
	w.script.Append(w.splitHorizontally())
	w.splitVerticallyRemaining()
	w.script.Append(selectNextPane())
}

// 'main-horizontal-flipped' layouts have one lower pane that is full width,
// and then split the remaining panes vertically along the top.
func (w *oldItermWindow) mainHorizontalFlipped() {
	w.script.Append(w.splitHorizontally())
	w.script.Append(selectPrevPane())
	w.splitVerticallyRemaining()
}

// 'tiled' layouts create 2 columns and then however many rows as
// needed. If there are odd number of panes then the bottom pane
// spans two columns. Panes are numbered top to bottom, left to right.
func (w *oldItermWindow) tiled() {
	for i := 2; i <= w.panes; i++ {
		if i%2 == 1 {
			w.script.Append(selectPrevPane(), w.splitHorizontally())
		} else {
			w.script.Append(w.splitVertically())
		}
	}
}

// '3_columns' layouts create 3 columns and then however many rows as
// needed. If there are odd number of panes then the bottom pane
// spans two columns. Panes are numbered top to bottom, left to right.
func (w *oldItermWindow) threeColumns() {
	for i := 2; i <= w.panes; i++ {
		if i%3 == 1 {
			w.script.Append(selectPrevPane(), selectPrevPane(), w.splitHorizontally())
		} else {
			w.script.Append(w.splitVertically())
		}
	}
}

// 'double-main-horizontal' layouts have two left panes that are full height,
// and then split the remaining panes horizontally down the right.
func (w *oldItermWindow) doubleMainHorizontal() {
	w.script.Append(w.splitVertically())
	if w.panes >= 2 {
		w.script.Append(w.splitVertically())
		w.splitHorizontallyRemaining()
	}
}

// 'double-main-vertical' layouts have two bottom panes that split the width
// and then split the remaining panes vertically across the top.
func (w *oldItermWindow) doubleMainVertical() {
	w.script.Append(w.splitHorizontally())
	w.script.Append(w.splitVertically())
	w.script.Append(selectNextPane())
	w.splitVerticallyRemaining()
}
