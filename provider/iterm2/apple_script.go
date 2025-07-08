package iterm2

import (
	"fmt"
	"strings"

	"github.com/boazj/muxocil/utils"
)

type AppleScriptKeyModifiers string

const (
	Command AppleScriptKeyModifiers = "command down"
	Option  AppleScriptKeyModifiers = "option down"
	Shift   AppleScriptKeyModifiers = "shift down"
)

type AsCmd string

type AppleScript struct {
	cmds []AsCmd
}

func newAppleScript() *AppleScript {
	return &AppleScript{
		cmds: make([]AsCmd, 0),
	}
}

func SingletonScript(cmd AsCmd) *AppleScript {
	return &AppleScript{
		cmds: []AsCmd{cmd},
	}
}

func (s *AppleScript) Append(cmds ...AsCmd) *AppleScript {
	s.cmds = append(s.cmds, cmds...)
	return s
}

func (s *AppleScript) Merge(script *AppleScript) *AppleScript {
	s.Append(script.cmds...)
	return s
}

func (s *AppleScript) Raw() string {
	return strings.Join(utils.ToStringSlice(s.cmds), "\n")
}

func (s *AppleScript) Execute() {
	// TODO:
}

func selectNextPane() AsCmd {
	return pressKeystroke("]", Command)
}

func selectPrevPane() AsCmd {
	return pressKeystroke("[", Command)
}

func pressKeystroke(key string, modifiers ...AppleScriptKeyModifiers) AsCmd {
	mods := ""
	if len(modifiers) == 1 {
		mods = string(modifiers[1])
	} else if len(modifiers) > 1 {
		mods = fmt.Sprintf("{%s}", strings.Join(utils.ToStringSlice(modifiers), ", "))
	}
	return AsCmd(fmt.Sprintf("tell i term application \"System Events\" to keystroke \"%s\" using %s", key, mods))
}

func selectColumnTopPane() AsCmd {
	return pressKeyCode("125", Command, Option)
}

func pressKeyCode(key string, modifiers ...AppleScriptKeyModifiers) AsCmd {
	mods := ""
	if len(modifiers) == 1 {
		mods = string(modifiers[1])
	} else if len(modifiers) > 1 {
		mods = fmt.Sprintf("{%s}", strings.Join(utils.ToStringSlice(modifiers), ", "))
	}
	return AsCmd(fmt.Sprintf("tell i term application \"System Events\" to key code \"%s\" using %s", key, mods))
}
