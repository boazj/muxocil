package tmux

import (
	"fmt"
	"os"
	"slices"
	"strconv"

	"github.com/boazj/muxocil/common"
	"github.com/charmbracelet/log"
)

type IdentifierKind string

const (
	idKind       IdentifierKind = "id"
	indexKind    IdentifierKind = "index"
	nameKind     IdentifierKind = "name"
	tokenKind    IdentifierKind = "token"
	scriptKind   IdentifierKind = "script"
	currentKind  IdentifierKind = "current"
	infferedKind IdentifierKind = "inffered"
)

var (
	WinAllowedTokens  = []string{"{start}", "{end}", "{last}", "{next}", "{previous}"}
	WinAllowedShort   = []string{"^", "$", "!"}
	WinAllowedPrefixs = []string{"+", "-"}
	PaneAllowedTokens = []string{
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
	PaneAllowedShort   = []string{"!"}
	PaneAllowedPrefixs = []string{"+", "-"}
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

func ID(val uint64) protoIdentifier {
	return protoIdentifier{
		value: strconv.FormatUint(val, 10),
		kind:  idKind,
	}
}

func Index(val uint64) protoIdentifier {
	return protoIdentifier{
		value: strconv.FormatUint(val, 10),
		kind:  indexKind,
	}
}

func Name(val string) protoIdentifier {
	return protoIdentifier{
		value: val,
		kind:  nameKind,
	}
}

func Token(val string) protoIdentifier {
	return protoIdentifier{
		value: val,
		kind:  tokenKind,
	}
}

func Script(val string) protoIdentifier {
	return protoIdentifier{
		value: val,
		kind:  scriptKind,
	}
}

func Current() protoIdentifier {
	return protoIdentifier{
		value: "",
		kind:  currentKind,
	}
}

func Inffered() protoIdentifier {
	return protoIdentifier{
		value: "",
		kind:  infferedKind,
	}
}

func (s *SessionIdentifier) String() string {
	switch s.kind {
	case idKind:
		if _, err := strconv.Atoi(s.value); err != nil {
			log.Fatal("session id should be a number, found ", "value", s.value)
			os.Exit(common.ExitProviderDataError)
		}
		return fmt.Sprintf("$%s", s.value)

	case nameKind:
		// Avoid partioal match using the = sign
		return fmt.Sprintf("=%s", s.value)

	case scriptKind:
		if s.value == "" {
			log.Fatal("session script value should be a valid string, found \"\"")
			os.Exit(common.ExitProviderDataError)
		}
		return s.value

	case currentKind:
		return ""
	case infferedKind:
		return ""
	}
	log.Fatal("session kind should be either of [id, name, script, current], found ", "kind", s.kind)
	os.Exit(common.ExitProviderDataError)
	return ""
}

func (w *WindowIdentifier) String() string {
	switch w.kind {
	case idKind:
		if _, err := strconv.Atoi(w.value); err != nil {
			log.Fatal("window id should be a number, found ", "value", w.value)
			os.Exit(common.ExitProviderDataError)
		}
		return fmt.Sprintf("@%s", w.value)

	case indexKind:
		if _, err := strconv.Atoi(w.value); err != nil {
			log.Fatal("window index should be a number, found ", "value", w.value)
			os.Exit(common.ExitProviderDataError)
		}
		return w.value

	case nameKind:
		// Avoid partioal match using the = sign
		return fmt.Sprintf("=%s", w.value)

	case scriptKind:
		if w.value == "" {
			log.Fatal("window script value should be a valid string, found \"\"")
			os.Exit(common.ExitProviderDataError)
		}
		return w.value

	case currentKind:
		return ""
	case infferedKind:
		return ""
	case tokenKind:
		if !isValidWindowToken(w.value) {
			log.Fatal("window token should be a legal value, found ", "value", w.value)
			os.Exit(common.ExitProviderDataError)
		}
		return w.value
	}
	log.Fatal("window kind should be either of [id, index, name, script, current], found ", "kind", w.kind)
	os.Exit(common.ExitProviderDataError)
	return ""
}

func (p *PaneIdentifier) String() string {
	switch p.kind {
	case idKind:
		if _, err := strconv.Atoi(p.value); err != nil {
			log.Fatal("pane id should be a number, found ", "value", p.value)
			os.Exit(common.ExitProviderDataError)
		}
		return fmt.Sprintf("%%%s", p.value)

	case indexKind:
		if _, err := strconv.Atoi(p.value); err != nil {
			log.Fatal("pane index should be a number, found ", "value", p.value)
			os.Exit(common.ExitProviderDataError)
		}
		return p.value

	case scriptKind:
		if p.value == "" {
			log.Fatal("pane script value should be a valid string, found \"\"")
			os.Exit(common.ExitProviderDataError)
		}
		return p.value

	case currentKind:
		return ""
	case tokenKind:
		if !isValidPaneToken(p.value) {
			log.Fatal("pane token should be a legal value, found ", "value", p.value)
			os.Exit(common.ExitProviderDataError)
		}
		return p.value
	}
	log.Fatal("pane kind should be either of [id, index, script, current], found ", "kind", p.kind)
	os.Exit(common.ExitProviderDataError)
	return ""
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

func TargetPaneFromWindow(window *TargetWindow, pane protoIdentifier) *TargetPane {
	return &TargetPane{
		session: window.session,
		window:  window.window,
		pane:    PaneIdentifier(pane),
	}
}

func (t *TargetSession) String() string {
	return t.session.String()
}

func (t *TargetWindow) String() string {
	if t.session.kind == currentKind {
		return t.window.String()
	}
	return fmt.Sprintf("%s:%s", t.session.String(), t.window.String())
}

func (t *TargetPane) String() string {
	if t.session.kind == currentKind && t.window.kind == currentKind {
		return t.pane.String()
	}
	if t.session.kind == currentKind {
		return fmt.Sprintf("%s.%s", t.window.String(), t.pane.String())
	}
	return fmt.Sprintf("%s:%s.%s", t.session.String(), t.window.String(), t.pane.String())
}

func (t *TargetSession) IsCurrent() bool {
	return t.session.kind == currentKind
}

func (t *TargetWindow) IsCurrent() bool {
	return t.session.kind == currentKind && t.window.kind == currentKind
}

func (t *TargetPane) IsCurrent() bool {
	return t.session.kind == currentKind && t.window.kind == currentKind && t.pane.kind == currentKind
}

func isValidWindowToken(value string) bool {
	if len(value) < 1 {
		return false
	}
	if len(value) == 1 && slices.Contains(WinAllowedShort, value) {
		return true
	}
	if slices.Contains(WinAllowedTokens, value) {
		return true
	}
	prefix := value[:1]
	if !slices.Contains(WinAllowedPrefixs, prefix) {
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
	if len(value) == 1 && slices.Contains(PaneAllowedShort, value) {
		return true
	}
	if slices.Contains(PaneAllowedTokens, value) {
		return true
	}
	prefix := value[:1]
	if !slices.Contains(PaneAllowedPrefixs, prefix) {
		return false
	}
	offset := value[1:]
	if _, err := strconv.Atoi(offset); err != nil {
		return false
	}
	return true
}
