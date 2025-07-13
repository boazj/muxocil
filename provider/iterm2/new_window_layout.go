package iterm2

import (
	"fmt"

	"github.com/boazj/muxocil/common"
)

type newItermWindow struct {
	panes        int
	createdPanes int

	layout common.MuxLayout

	layoutStrategy map[common.MuxLayout]func()

	script AppleScript
}

func newNewItermWindow(panes int, layout common.MuxLayout) *newItermWindow {
	if layout == "" {
		layout = common.Tiled
	}

	// tmux seems to treat the first 2 tiles of a tiled layout like this
	if panes == 2 && layout == common.Tiled {
		layout = common.EvenVertical
	}

	win := &newItermWindow{
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

func (w *newItermWindow) arrangePanes() (*AppleScript, error) {
	w.script.Append(AsCmd("set pane_1 to (current session of current window)"))

	if w.panes <= 1 {
		// If we have just one pane we don't need to do any splitting.
		return &AppleScript{}, nil
	}
	strat, ok := w.layoutStrategy[w.layout]
	if !ok {
		return &AppleScript{}, fmt.Errorf("unknown layout setting")
	}
	strat()

	return &w.script, nil
}

func (w *newItermWindow) createPane(parentPaneIndex int, childPaneIndex int, split string) AsCmd {
	if split == "" {
		split = "vertical"
	}
	return Aprintf(` tell pane_%d
     set pane_%d to (split %sly with same profile)
 end tell`, parentPaneIndex, childPaneIndex, split)
}

// 'even-horizontal' layouts just split vertically across the screen
func (w *newItermWindow) evenHorizontal() {
	for i := 2; i < w.panes+1; i++ {
		w.script.Append(w.createPane(i-1, i, "vertical"))
	}
}

// 'even-vertical' layouts just split horizontally down the screen
func (w *newItermWindow) evenVertical() {
	for i := 2; i < w.panes+1; i++ {
		w.script.Append(w.createPane(i-1, i, "horizontal"))
	}
}

// 'main-vertical' layouts have one left pane that is full height,
// and then split the remaining panes horizontally down the right
func (w *newItermWindow) mainVertical() {
	w.script.Append(w.createPane(1, 2, "vertical"))
	for i := 3; i < w.panes+1; i++ {
		w.script.Append(w.createPane(i-1, i, "horizontal"))
	}
}

// 'main-vertical-flipped' layouts have one right pane that is full height,
// and then split the remaining panes horizontally down the left
func (w *newItermWindow) mainVerticalFlipped() {
	w.script.Append(w.createPane(1, w.panes, "vertical"))
	for i := 2; i < w.panes; i++ {
		w.script.Append(w.createPane(i-1, i, "horizontal"))
	}
}

// 'main-horizontal' layouts have one upper pane that is full width,
// and then split the remaining panes vertically along the buttom
func (w *newItermWindow) mainHorizontal() {
	w.script.Append(w.createPane(1, 2, "horizontal"))
	for i := 3; i < w.panes+1; i++ {
		w.script.Append(w.createPane(i-1, i, "vertical"))
	}
}

// 'main-horizontal-flipped' layouts have one lower pane that is full width,
// and then split the remaining panes vertically along the top
func (w *newItermWindow) mainHorizontalFlipped() {
	w.script.Append(w.createPane(1, w.panes, "horizontal"))
	for i := 2; i < w.panes; i++ {
		w.script.Append(w.createPane(i-1, i, "vertical"))
	}
}

// 'tiled' layouts create 2 columns and then however many rows as
// needed. If there are odd number of panes then the bottom pane
// spans two columns. Panes are numbered top to bottom, left to right.
func (w *newItermWindow) tiled() {
	for i := 2; i <= w.panes; i++ {
		if i%2 == 1 {
			w.script.Append(w.createPane(i-2, i, "horizontal"))
		} else {
			w.script.Append(w.createPane(i-1, i, "vertical"))
		}
	}
}

// '3_columns' layouts create 3 columns and then however many rows as
// needed. If there are odd number of panes then the bottom pane
// spans two columns. Panes are numbered top to bottom, left to right.
func (w *newItermWindow) threeColumns() {
	for i := 2; i <= w.panes; i++ {
		if i%3 == 1 {
			w.script.Append(w.createPane(i-3, i, "horizontal"))
		} else {
			w.script.Append(w.createPane(i-1, i, "vertical"))
		}
	}
}

// 'double-main-horizontal' layouts have two bottom panes that split the width
// and then split the remaining panes vertically across the top
func (w *newItermWindow) doubleMainHorizontal() {
	w.script.Append(w.createPane(1, w.panes-1, "horizontal"))
	for i := 2; i < w.panes-1; i++ {
		w.script.Append(w.createPane(i-1, i, "vertical"))
	}
}

// 'double-main-vertical' layouts have two left panes that are full height,
// and then split the remaining panes horizontally down the right
func (w *newItermWindow) doubleMainVertical() {
	w.script.Append(w.createPane(1, 2, "vertical"))
	w.script.Append(w.createPane(2, 3, "vertical"))
	for i := 4; i < w.panes+1; i++ {
		w.script.Append(w.createPane(i-1, i, "horizontal"))
	}
}
