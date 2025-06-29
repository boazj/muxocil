package tmux

import (
	"fmt"
	"slices"
	"strconv"
)

type IdentifierKind string

const (
	ID       IdentifierKind = "id"
	INDEX    IdentifierKind = "index"
	NAME     IdentifierKind = "name"
	TOKEN    IdentifierKind = "token"
	SCRIPT   IdentifierKind = "script"
	CURRENT  IdentifierKind = "current"
	INFFERED IdentifierKind = "inffered"
)

var (
	WINDOW_ALLOWED_TOKENS  = []string{"{start}", "{end}", "{last}", "{next}", "{previous}"}
	WINDOW_ALLOWED_SHORT   = []string{"^", "$", "!"}
	WINDOW_ALLOWED_PREFIXS = []string{"+", "-"}
	PANE_ALLOWED_TOKENS    = []string{
		"{last}",
		"{next}",
		"{previous}",
		"{top}",
		"{bottom}",
		"{left}",
		"{right}",
		"{top-left}",
		"{top-right}",
		"{bottom-left}",
		"{bottom-right}",
		"{up-of}",
		"{down-of}",
		"{left-of}",
		"{right-of}",
	}
	PANE_ALLOWED_SHORT   = []string{"!"}
	PANE_ALLOWED_PREFIXS = []string{"+", "-"}
)

type identifier struct {
	value string
	kind  IdentifierKind
}

type (
	protoIdentifier identifier

	SessionIdentifier identifier
	WindowIdentifier  identifier
	PaneIdentifier    identifier
)

func Id(val uint64) protoIdentifier {
	return protoIdentifier{
		value: strconv.FormatUint(val, 10),
		kind:  ID,
	}
}

func Index(val uint64) protoIdentifier {
	return protoIdentifier{
		value: strconv.FormatUint(val, 10),
		kind:  INDEX,
	}
}

func Name(val string) protoIdentifier {
	return protoIdentifier{
		value: val,
		kind:  NAME,
	}
}

func Token(val string) protoIdentifier {
	return protoIdentifier{
		value: val,
		kind:  TOKEN,
	}
}

func Script(val string) protoIdentifier {
	return protoIdentifier{
		value: val,
		kind:  SCRIPT,
	}
}

func Current() protoIdentifier {
	return protoIdentifier{
		value: "",
		kind:  CURRENT,
	}
}

func Inffered() protoIdentifier {
	return protoIdentifier{
		value: "",
		kind:  INFFERED,
	}
}

func (s *SessionIdentifier) String() string {
	switch s.kind {
	case ID:
		if _, err := strconv.Atoi(s.value); err != nil {
			panic(fmt.Sprintf("session id should be a number, found [%s]", s.value))
		}
		return fmt.Sprintf("$%s", s.value)

	case NAME:
		// Avoid partioal match using the = sign
		return fmt.Sprintf("=%s", s.value)

	case SCRIPT:
		if s.value == "" {
			panic("session script value should be a valid string, found \"\"")
		}
		return s.value

	case CURRENT:
		return ""
	case INFFERED:
		return ""
	}
	panic(fmt.Sprintf("session kind should be either of [id, name, script, current], found [%s]", s.kind))
}

func (w *WindowIdentifier) String() string {
	switch w.kind {
	case ID:
		if _, err := strconv.Atoi(w.value); err != nil {
			panic(fmt.Sprintf("window id should be a number, found [%s]", w.value))
		}
		return fmt.Sprintf("@%s", w.value)

	case INDEX:
		if _, err := strconv.Atoi(w.value); err != nil {
			panic(fmt.Sprintf("window index should be a number, found [%s]", w.value))
		}
		return w.value

	case NAME:
		// Avoid partioal match using the = sign
		return fmt.Sprintf("=%s", w.value)

	case SCRIPT:
		if w.value == "" {
			panic("window script value should be a valid string, found \"\"")
		}
		return w.value

	case CURRENT:
		return ""
	case INFFERED:
		return ""
	case TOKEN:
		if !isValidWindowToken(w.value) {
			panic(fmt.Sprintf("window token should be a legal value, found [%s]", w.value))
		}
		return w.value
	}
	panic(fmt.Sprintf("window kind should be either of [id, index, name, script, current], found [%s]", w.kind))
}

func (p *PaneIdentifier) String() string {
	switch p.kind {
	case ID:
		if _, err := strconv.Atoi(p.value); err != nil {
			panic(fmt.Sprintf("pane id should be a number, found [%s]", p.value))
		}
		return fmt.Sprintf("%%%s", p.value)

	case INDEX:
		if _, err := strconv.Atoi(p.value); err != nil {
			panic(fmt.Sprintf("pane index should be a number, found [%s]", p.value))
		}
		return p.value

	case SCRIPT:
		if p.value == "" {
			panic("pane script value should be a valid string, found \"\"")
		}
		return p.value

	case CURRENT:
		return ""
	case TOKEN:
		if !isValidPaneToken(p.value) {
			panic(fmt.Sprintf("pane token should be a legal value, found [%s]", p.value))
		}
		return p.value
	}
	panic(fmt.Sprintf("pane kind should be either of [id, index, script, current], found [%s]", p.kind))
}

type target struct {
	session SessionIdentifier
	window  WindowIdentifier
	pane    PaneIdentifier
}

type (
	TargetSession target
	TargetWindow  target
	TargetPane    target
)

func NewTargetSession(session protoIdentifier) TargetSession {
	return TargetSession{
		session: SessionIdentifier(session),
	}
}

func NewTargetWindow(session protoIdentifier, window protoIdentifier) TargetWindow {
	return TargetWindow{
		session: SessionIdentifier(session),
		window:  WindowIdentifier(window),
	}
}

func NewTargetPane(session protoIdentifier, window protoIdentifier, pane protoIdentifier) TargetPane {
	return TargetPane{
		session: SessionIdentifier(session),
		window:  WindowIdentifier(window),
		pane:    PaneIdentifier(pane),
	}
}

func TargetPaneFromWindow(window TargetWindow, pane protoIdentifier) TargetPane {
	return TargetPane{
		session: window.session,
		window:  window.window,
		pane:    PaneIdentifier(pane),
	}
}

func (t *TargetSession) String() string {
	return t.session.String()
}

func (t *TargetWindow) String() string {
	if t.session.kind == CURRENT {
		return t.window.String()
	}
	return fmt.Sprintf("%s:%s", t.session.String(), t.window.String())
}

func (t *TargetPane) String() string {
	if t.session.kind == CURRENT && t.window.kind == CURRENT {
		return t.pane.String()
	}
	if t.session.kind == CURRENT {
		return fmt.Sprintf("%s.%s", t.window.String(), t.pane.String())
	}
	return fmt.Sprintf("%s:%s.%s", t.session.String(), t.window.String(), t.pane.String())
}

func (t *TargetSession) IsCurrent() bool {
	return t.session.kind == CURRENT
}

func (t *TargetWindow) IsCurrent() bool {
	return t.session.kind == CURRENT && t.window.kind == CURRENT
}

func (t *TargetPane) IsCurrent() bool {
	return t.session.kind == CURRENT && t.window.kind == CURRENT && t.pane.kind == CURRENT
}

func isValidWindowToken(value string) bool {
	if len(value) < 1 {
		return false
	}
	if len(value) == 1 && slices.Contains(WINDOW_ALLOWED_SHORT, value) {
		return true
	}
	if slices.Contains(WINDOW_ALLOWED_TOKENS, value) {
		return true
	}
	prefix := value[:1]
	if !slices.Contains(WINDOW_ALLOWED_PREFIXS, prefix) {
		return false
	}
	offset := value[1:]
	if _, err := strconv.Atoi(offset); err != nil {
		return false
	}
	return true
}

func isValidPaneToken(value string) bool {
	if len(value) < 1 {
		return false
	}
	if len(value) == 1 && slices.Contains(PANE_ALLOWED_SHORT, value) {
		return true
	}
	if slices.Contains(PANE_ALLOWED_TOKENS, value) {
		return true
	}
	prefix := value[:1]
	if !slices.Contains(PANE_ALLOWED_PREFIXS, prefix) {
		return false
	}
	offset := value[1:]
	if _, err := strconv.Atoi(offset); err != nil {
		return false
	}
	return true
}
