package iterm2

import (
	"fmt"
	"strings"

	"github.com/boazj/muxocil/utils"
)

//FIXME:
//lint:file-ignore U1000 in dev

type AppleScriptKeyModifiers string

const (
	KeyCommand AppleScriptKeyModifiers = "command down"
	KeyOption  AppleScriptKeyModifiers = "option down"
	KeyShift   AppleScriptKeyModifiers = "shift down"
)

type AsCmd string

type AppleScript struct {
	cmds       []AsCmd
	suffixCmds []AsCmd
}

func newAppleScript() *AppleScript {
	return &AppleScript{
		cmds:       make([]AsCmd, 0),
		suffixCmds: make([]AsCmd, 0),
	}
}

func SingletonScript(cmd AsCmd) *AppleScript {
	return newAppleScript().Append(cmd)
}

func Aprintf(format string, a ...any) AsCmd {
	return AsCmd(fmt.Sprintf(format, a...))
}

func (s *AppleScript) Append(cmds ...AsCmd) *AppleScript {
	s.cmds = append(s.cmds, cmds...)
	return s
}

func (s *AppleScript) Suffix(cmds ...AsCmd) *AppleScript {
	s.suffixCmds = append(s.suffixCmds, cmds...)
	return s
}

func (s *AppleScript) Merge(script *AppleScript) *AppleScript {
	s.Append(script.cmds...)
	s.Suffix(script.suffixCmds...)
	return s
}

func (s *AppleScript) Raw() string {
	main := strings.Join(utils.ToStringSlice(s.cmds), "\n")
	suffix := strings.Join(utils.ToStringSlice(s.suffixCmds), "\n")
	return main + "\n" + suffix
}

func (s *AppleScript) Execute() {
	// TODO:
}

func selectNextPane() AsCmd {
	return pressKeystroke("]", KeyCommand)
}

func selectPrevPane() AsCmd {
	return pressKeystroke("[", KeyCommand)
}

func pressKeystroke(key string, modifiers ...AppleScriptKeyModifiers) AsCmd {
	mods := ""
	if len(modifiers) == 1 {
		mods = string(modifiers[0])
	} else if len(modifiers) > 1 {
		mods = fmt.Sprintf("{%s}", strings.Join(utils.ToStringSlice(modifiers), ", "))
	}
	return Aprintf("tell %s \"System Events\" to keystroke \"%s\" using %s", TargetOldApp, key, mods)
}

func selectColumnTopPane() AsCmd {
	return pressKeyCode("125", KeyCommand, KeyOption)
}

func pressKeyCode(key string, modifiers ...AppleScriptKeyModifiers) AsCmd {
	mods := ""
	if len(modifiers) == 1 {
		mods = string(modifiers[1])
	} else if len(modifiers) > 1 {
		mods = fmt.Sprintf("{%s}", strings.Join(utils.ToStringSlice(modifiers), ", "))
	}
	return Aprintf("tell %s \"System Events\" to key code \"%s\" using %s", TargetOldApp, key, mods)
}
