package tmux

import (
	"strings"

	"github.com/boazj/muxocil/common"
)

type WindowLayoutStrategy interface {
	PlacePane(tmux *Tmux, window *common.Window, pane *common.Pane, index int) error
}

var strategies = map[string]WindowLayoutStrategy{
	"even-horizontal":          &nativeLayout{layout: "even-horizontal"},
	"even-vertical":            &nativeLayout{layout: "even-vertical"},
	"main-horizontal":          &nativeLayout{layout: "main-horizontal"},
	"main-horizontal-mirrored": &nativeLayout{layout: "main-horizontal-mirrored"},
	"main-vertical":            &nativeLayout{layout: "main-vertical"},
	"main-vertical-mirrored":   &nativeLayout{layout: "main-vertical-mirrored"},
	"main-vertical-flipped":    &nativeLayout{layout: "main-vertical-mirrored"},
	"tiled":                    &nativeLayout{layout: "tiled"},
	"3_columns":                &ThreeColumns{},
	"double-main-horizontal":   &DoubleMainHorizontal{},
	"double-main-vertical":     &DoubleMainVertical{},
}

// https://github.com/TomAnthony/itermocil/blob/master/LAYOUTS.md
type nativeLayout struct {
	layout string
}

func (l *nativeLayout) PlacePane(tmux *Tmux, window *common.Window, pane *common.Pane, index int) error {
	tWin := tmux.wins[window]
	var cmd string
	var tPane TargetPane
	// pane 0 already exists and configure by the window creation
	if index != 0 {
		tPane, cmd = tmux.splitWindow(
			TargetPaneFromWindow(tWin, Current()),
			window.Root,
			pane.Focus,
			true,
		)
		tmux.cmds = append(tmux.cmds, cmd)
	} else {
		tPane = TargetPaneFromWindow(tWin, Current())
	}
	if len(pane.Commands) > 0 {
		tmux.cmds = append(tmux.cmds, tmux.sendKeysToPane(tPane, strings.Join(pane.Commands, "; "), true))
	}
	if len(window.Panes)-1 == index {
		// last pane, so apply native layout
		tmux.cmds = append(tmux.cmds, tmux.selectLayout(tWin, l.layout))
	}
	return nil
}

// iterm - https://github.com/TomAnthony/itermocil/blob/master/LAYOUTS.md
//
//	3_columns - Creates 3 columns and then however many rows as needed. If the number of panes isn't divisible by 3 then the final row will have fewer columns.
//
// .------------.------------.------------.
// | (0)        | (1)        | (2)        |
// |            |            |            |
// |            |            |            |
// |------------|------------|------------|
// | (3)        | (4)        | (5)        |
// |            |            |            |
// |            |            |            |
// |------------|------------|------------|
// | (6)        | (7)        | (8)        |
// |            |            |            |
// |            |            |            |
// '------------'------------'------------'

type ThreeColumns struct{}

func (l *ThreeColumns) PlacePane(tmux *Tmux, window *common.Window, pane *common.Pane, index int) error {
	// TODO:
	return nil
}

// iterm - https://github.com/TomAnthony/itermocil/blob/master/LAYOUTS.md
//
//	double-main-horizontal - Create 2 rows. The bottom row is 2 full width columns and the top row is split into as many columns as needed.
//
// .-----------.-------------.-----------.
// | (0)       | (1)         | (2)       |
// |           |             |           |
// |           |             |           |
// |           |             |           |
// |-----------'------.------'-----------|
// | (4)              | (3)              |
// |                  |                  |
// |                  |                  |
// |                  |                  |
// |                  |                  |
// '------------------'------------------'

type DoubleMainHorizontal struct{}

func (l *DoubleMainHorizontal) PlacePane(tmux *Tmux, window *common.Window, pane *common.Pane, index int) error {
	// TODO:
	return nil
}

// iterm - https://github.com/TomAnthony/itermocil/blob/master/LAYOUTS.md
//
//	double-main-vertical - Create 2 full height columns on the left, and a third column with as many rows as needed.
//
// .-----------.-------------.-----------.
// | (0)       | (1)         | (2)       |
// |           |             |           |
// |           |             |-----------|
// |           |             | (3)       |
// |           |             |           |
// |           |             |-----------|
// |           |             | (4)       |
// |           |             |           |
// |           |             |           |
// '-----------'-------------'-----------'

type DoubleMainVertical struct{}

func (l *DoubleMainVertical) PlacePane(tmux *Tmux, window *common.Window, pane *common.Pane, index int) error {
	// TODO:
	return nil
}
