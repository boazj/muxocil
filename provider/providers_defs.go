package provider

import (
	"runtime"

	"github.com/boazj/muxocil/common"
	"github.com/boazj/muxocil/provider/iterm2"
	"github.com/boazj/muxocil/provider/kitty"
	"github.com/boazj/muxocil/provider/tmux"
	"github.com/boazj/muxocil/provider/wezterm"
	"github.com/boazj/muxocil/provider/zellij"
	"github.com/boazj/muxocil/utils"
	"github.com/tiendc/gofn"
)

type providerDefs struct {
	muxers map[common.MuxID]Mux
}

func CreateProviderDefs() *providerDefs {
	p := &providerDefs{
		muxers: make(map[common.MuxID]Mux),
	}
	allos := []common.OS{common.Windows, common.MacOS, common.Linux}
	unix := []common.OS{common.MacOS, common.Linux}
	mac := []common.OS{common.MacOS}
	p.pushMux(common.Tmux, "tmux", allos, tmux.NewProvider, tmux.Detect)
	p.pushMux(common.Zellij, "Zellij", unix, zellij.NewProvider, zellij.Detect)
	p.pushEmu(common.Wezterm, "WezTerm", allos, wezterm.NewProvider, wezterm.Detect)
	p.pushEmu(common.Iterm2, "iTerm2", mac, iterm2.NewProvider, iterm2.Detect)
	p.pushEmu(common.Kitty, "Kitty", unix, kitty.NewProvider, kitty.Detect)
	return p
}

func (p *providerDefs) GetSupportedProviders() []Mux {
	var curOs common.OS
	switch runtime.GOOS {
	case "darwin":
		curOs = common.MacOS
	case "windows":
		curOs = common.Windows
	case "linux":
		curOs = common.Linux
	default:
		// TODO: exit
	}
	return utils.MapFilterValues(p.muxers, func(v Mux) bool {
		return gofn.Contain(v.SupportedOs, curOs)
	})
}

func (p *providerDefs) pushEmu(
	id common.MuxID,
	display string,
	os []common.OS,
	construct func(*common.CommandOpts) (common.Provider, error),
	detector func(*common.Config) (bool, bool),
) {
	p.push(id, display, common.Emulator, os, construct, detector)
}

func (p *providerDefs) pushMux(
	id common.MuxID,
	display string,
	os []common.OS,
	construct func(*common.CommandOpts) (common.Provider, error),
	detector func(*common.Config) (bool, bool),
) {
	p.push(id, display, common.Multiplexer, os, construct, detector)
}

func (p *providerDefs) push(
	id common.MuxID,
	display string,
	kind common.ProviderType,
	os []common.OS,
	construct func(*common.CommandOpts) (common.Provider, error),
	detector func(*common.Config) (bool, bool),
) {
	p.muxers[id] = Mux{
		id,
		display,
		kind,
		os,
		construct,
		detector,
	}
}

func (p *providerDefs) GetSupportedEmulators() []Mux {
	return utils.MapFilterValues(p.muxers, func(v Mux) bool {
		return v.Kind == common.Emulator
	})
}

func (p *providerDefs) GetSupportedMultiplexers() []Mux {
	return utils.MapFilterValues(p.muxers, func(v Mux) bool {
		return v.Kind == common.Multiplexer
	})
}

func (p *providerDefs) GetAllProviders() []Mux {
	return gofn.MapValues(p.muxers)
}

func (p *providerDefs) GetProvider(id common.MuxID) (Mux, bool) {
	m, ok := p.muxers[id]
	return m, ok
}

type Mux struct {
	ID          common.MuxID
	Display     string
	Kind        common.ProviderType
	SupportedOs []common.OS
	constructor func(*common.CommandOpts) (common.Provider, error)
	detector    func(*common.Config) (bool, bool) // is in app, can run app
}
