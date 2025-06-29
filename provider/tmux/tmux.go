// Package tmux represents the provider for the tmux Multiplexr
package tmux

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/boazj/muxocil/common"
)

type Tmux struct {
	opts *common.CommandOpts

	baseIndex     int
	paneBaseIndex int

	// Not elegant: some stateful crap to optimize tmux session-window creation
	createSession bool
	sessionName   string

	wins map[*common.Window]TargetWindow

	cmds []string
}

func NewTmux(opts *common.CommandOpts) (*Tmux, error) {
	baseIndex, err := getTmuxOptionNumericValue("base-index")
	if err != nil {
		return nil, err
	}

	paneBaseIndex, err := getTmuxOptionNumericValue("pane-base-index")
	if err != nil {
		return nil, err
	}

	t := Tmux{
		opts:          opts,
		cmds:          make([]string, 0),
		wins:          make(map[*common.Window]TargetWindow),
		createSession: false,
		baseIndex:     baseIndex,
		paneBaseIndex: paneBaseIndex,
	}

	return &t, nil
}

func (t *Tmux) CreateLayout(layout *common.Layout) error {
	if layout.Name == "" {
		return nil
	}
	t.createSession = true
	t.sessionName = layout.Name
	return nil
}

func (t *Tmux) CreateWindow(window *common.Window, index int) error {
	var tWin TargetWindow
	var cmd string
	if index == 0 && t.createSession {
		tWin, cmd = t.newSession(t.sessionName, window.Name, window.Root, window.Focus)
	} else {
		tWin, cmd = t.newWindow(window.Name, window.Root, window.Focus)
	}
	t.cmds = append(t.cmds, cmd)
	t.wins[window] = tWin

	for k, v := range window.Options {
		t.cmds = append(t.cmds, t.setWindowOption(tWin, k, v))
	}
	t.cmds = append(t.cmds, t.selectLayout(tWin, window.Layout))
	if len(window.Commands) > 0 {
		t.cmds = append(t.cmds, t.sendKeysToWindow(tWin, strings.Join(window.Commands, "; "), true))
	}
	return nil
}

func (t *Tmux) CreatePane(window *common.Window, pane *common.Pane, index int) error {
	// TODO:
	return nil
}

func (t *Tmux) GetCommads() []string {
	return t.cmds
}

func getTmuxOptionValue(option string) (string, error) {
	cmd := exec.Command(CMD, "show-option", "-gv", option)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("could not fetch tmux %s option value: %v", option, err)
	}

	return strings.TrimSpace(string(out)), nil
}

func getTmuxOptionNumericValue(option string) (int, error) {
	raw, err := getTmuxOptionValue(option)
	if err != nil {
		return -1, err
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return -1, fmt.Errorf("could not parse tmux %s option value to int: %v", option, err)
	}
	return value, nil
}
