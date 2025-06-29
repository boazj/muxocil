package tmux

import "fmt"

const CMD = "tmux"

const (
	WinID        = "winid"
	PaneID       = "paneid"
	SessionWinID = "sessionwinid"
)

func (t *Tmux) newSession(sessionName string, windowName string, rootDirectory string, focus bool) (TargetWindow, string) {
	sessionNameFlag := ""
	if sessionName != "" {
		sessionNameFlag = fmt.Sprintf("-s '%s'", sessionName)
	}

	winNameFlag := ""
	if windowName != "" {
		winNameFlag = fmt.Sprintf("-n '%s'", windowName)
	}

	root := ""
	if rootDirectory != "" {
		root = fmt.Sprintf("-c '%s'", rootDirectory) // TODO: consider \" instead of \' to support vars in root
	}

	focusFlag := "-d"
	if focus {
		focusFlag = ""
	}
	cmd := fmt.Sprintf(
		"%s=\"$(%s new-session %s %s %s %s -P -F \"#{session_id}:#{window_id}\")\"",
		SessionWinID,
		CMD,
		focusFlag,
		sessionNameFlag,
		winNameFlag,
		root,
	)
	target := NewTargetWindow(Inffered(), Script(fmt.Sprintf("$%s", SessionWinID)))
	return target, cmd
}

func (t *Tmux) newWindow(newName string, rootDirectory string, focus bool) (TargetWindow, string) {
	name := ""
	if newName != "" {
		name = fmt.Sprintf("-n '%s'", newName)
	}
	root := ""
	if rootDirectory != "" {
		root = fmt.Sprintf("-c '%s'", rootDirectory) // TODO: consider \" instead of \' to support vars in root
	}
	focusFlag := "-d"
	if focus {
		focusFlag = ""
	}
	cmd := fmt.Sprintf(
		"%s=\"$(%s new-window %s %s %s -P -F \"#{session_id}:#{window_id}\")\"",
		WinID,
		CMD,
		focusFlag,
		name,
		root,
	)
	target := NewTargetWindow(Current(), Script(fmt.Sprintf("$%s", WinID)))
	return target, cmd
}

func (t *Tmux) renameWindow(window TargetWindow, newName string) string {
	target := ""
	if !window.IsCurrent() {
		target = fmt.Sprintf("-t %s", window.String())
	}
	return fmt.Sprintf("%s rename-window %s '%s'", CMD, target, newName)
}

func (t *Tmux) setWindowOption(window TargetWindow, option string, value string) string {
	target := ""
	if !window.IsCurrent() {
		target = fmt.Sprintf("-t %s", window.String())
	}
	return fmt.Sprintf("%s set-window-option %s %s %s", CMD, target, option, value)
}

func (t *Tmux) listPanes(window TargetWindow) string {
	target := ""
	if !window.IsCurrent() {
		target = fmt.Sprintf("-t %s", window.String())
	}
	return fmt.Sprintf("%s list-panes %s", CMD, target)
}

func (t *Tmux) listWindows() string {
	return "%s list-windows"
}

func (t *Tmux) renameSession(session TargetSession, newName string) string {
	target := ""
	if !session.IsCurrent() {
		target = fmt.Sprintf("-t %s", session.String())
	}
	return fmt.Sprintf("%s rename-session %s '%s'", CMD, target, newName)
}

func (t *Tmux) selectLayout(window TargetWindow, layout string) string {
	if layout == "" {
		return ""
	}
	target := ""
	if !window.IsCurrent() {
		target = fmt.Sprintf("-t %s", window.String())
	}
	return fmt.Sprintf("%s select-layout %s '%s'", CMD, target, layout)
}

func (t *Tmux) selectWindow(window TargetWindow) string {
	if window.IsCurrent() {
		return "" // equivalent to selecting the selected window
	}
	return fmt.Sprintf("%s select-window %s", CMD, window.String())
}

func (t *Tmux) selectPane(pane TargetPane) string {
	if pane.IsCurrent() {
		return "" // equivalent to selecting the selected pane
	}
	return fmt.Sprintf("%s select-pane %s", CMD, pane.String())
}

func (t *Tmux) sendKeys(keys string) string {
	return fmt.Sprintf("%s send-keys '%s'", CMD, keys)
}

func (t *Tmux) sendKeysToWindow(window TargetWindow, keys string, enter bool) string {
	targetPane := TargetPaneFromWindow(window, Token("{last}"))
	return t.sendKeysToPane(targetPane, keys, enter)
}

func (t *Tmux) sendKeysToPane(pane TargetPane, keys string, enter bool) string {
	target := ""
	if !pane.IsCurrent() {
		target = fmt.Sprintf("-t %s", pane.String())
	}
	execute := ""
	if enter {
		execute = fmt.Sprintf("\\; send-keys %s Enter", target)
	}
	return fmt.Sprintf("%s send-keys %s -l '%s' %s", CMD, target, keys, execute)
}

func (t *Tmux) showOptions(name string) string {
	return fmt.Sprintf("%s show-options -gv %s", CMD, name)
}

func (t *Tmux) showWindowOptions(name string) string {
	return fmt.Sprintf("%s show-window-options -gv %s", CMD, name)
}

func (t *Tmux) splitWindow(pane TargetPane, rootDirecotry string, vertical bool, focus bool) (TargetPane, string) {
	srcTarget := ""
	if !pane.IsCurrent() {
		srcTarget = fmt.Sprintf("-t %s", pane.String())
	}
	root := ""
	if rootDirecotry != "" {
		root = fmt.Sprintf("-c '%s'", rootDirecotry)
	}
	direction := "-v"
	if !vertical {
		direction = "-h"
	}
	focusFlag := "-d"
	if focus {
		focusFlag = ""
	}
	cmd := fmt.Sprintf(
		"%s=\"$(%s split-window %s %s %s %s -P -F \"#{session_id}:#{window_id}.#{pane_id})\"",
		PaneID,
		CMD,
		focusFlag,
		direction,
		srcTarget,
		root,
	)
	target := NewTargetPane(Inffered(), Inffered(), Script(fmt.Sprintf("$%s", PaneID)))
	return target, cmd
}
