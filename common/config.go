// Package common contains common types for muxocil
package common

import (
	"os"

	"github.com/boazj/muxocil/utils"
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
	"github.com/tiendc/gofn"
)

const (
	ConfSearchLocations               = "configuration.layout_search_locations"
	ConfSearchExcludeLocations        = "configuration.layout_exclude_locations"
	ProviderUseProvider               = "providers.use_provider"
	ProviderUseProviderFail           = "providers.use_provider_fail_fast"
	ProviderValidateProvider          = "providers.validate_provider_from_env"
	ProviderInfferProvider            = "providers.inffer_provider_from_env"
	ProviderLaunchMux                 = "providers.launch_multiplexer_app"
	LayoutPaneShowTitle               = "layout.panes.show_pane_title"
	LayoutPaneShowTitleWhitespace     = "layout.panes.pane_title_replace_whitespace"
	LayoutPaneShowTitleWhitespaceChar = "layout.panes.pane_title_replace_whitespace_char"
	LayoutWinIgnoreLayout             = "layout.window.ignore_window_layout"
	LayoutSessionRepalceIfExists      = "layout.session.replace_if_session_exists"
)

type Config struct {
	LayoutSearchLocations          []string
	LayoutSearchExcludeLocations   []string
	UseProvider                    string
	UseProviderFailFast            bool
	ValidateProviderFromEnv        bool
	InfferProviderFromEnv          bool
	LaunchMultiplexerApp           bool
	ShowPaneTitle                  bool
	PaneTitleReplaceWhitespace     bool
	PaneTitleReplaceWhitespaceChar string
	IgnoreWindowLayout             bool
	ReplaceIfSessionExists         bool
}

func GetConfig() *Config {
	conf := Config{}

	viper.SetDefault(ConfSearchLocations, []string{"$HOME/.teamocil"})
	viper.SetDefault(ConfSearchExcludeLocations, []string{})
	viper.SetDefault(ProviderUseProvider, "")
	viper.SetDefault(ProviderUseProviderFail, true)
	viper.SetDefault(ProviderValidateProvider, true)
	viper.SetDefault(ProviderInfferProvider, true)
	viper.SetDefault(ProviderLaunchMux, false)
	viper.SetDefault(LayoutPaneShowTitle, true)
	viper.SetDefault(LayoutPaneShowTitleWhitespace, false)
	viper.SetDefault(LayoutPaneShowTitleWhitespaceChar, "_")
	viper.SetDefault(LayoutWinIgnoreLayout, false)
	viper.SetDefault(LayoutSessionRepalceIfExists, true)
	conf.LayoutSearchLocations = viper.GetStringSlice(ConfSearchLocations)
	conf.LayoutSearchExcludeLocations = viper.GetStringSlice(ConfSearchExcludeLocations)
	conf.UseProvider = viper.GetString(ProviderUseProvider)
	conf.UseProviderFailFast = viper.GetBool(ProviderUseProviderFail)
	conf.ValidateProviderFromEnv = viper.GetBool(ProviderValidateProvider)
	conf.InfferProviderFromEnv = viper.GetBool(ProviderInfferProvider)
	conf.LaunchMultiplexerApp = viper.GetBool(ProviderLaunchMux)
	conf.ShowPaneTitle = viper.GetBool(LayoutPaneShowTitle)
	conf.PaneTitleReplaceWhitespace = viper.GetBool(LayoutPaneShowTitleWhitespace)
	conf.PaneTitleReplaceWhitespaceChar = viper.GetString(LayoutPaneShowTitleWhitespaceChar)
	conf.IgnoreWindowLayout = viper.GetBool(LayoutWinIgnoreLayout)
	conf.ReplaceIfSessionExists = viper.GetBool(LayoutSessionRepalceIfExists)

	return &conf
}

func (c *Config) GetLayoutSearchLocations() []string {
	log.Debug("Reading search locations", "locations", c.LayoutSearchLocations)
	elocs := gofn.ToSet(gofn.MapSlice(c.LayoutSearchLocations, os.ExpandEnv))
	log.Debug("Expanding search locations", "expanded_locations", elocs)

	log.Debug("Reading exclude locations", "locations", c.LayoutSearchExcludeLocations)
	excludes := gofn.ToSet(gofn.MapSlice(c.LayoutSearchExcludeLocations, os.ExpandEnv))
	log.Debug("Expanding exclude locations", "locations", excludes)

	elocs = gofn.FilterNIN(elocs, excludes...)
	log.Debug("Effective search locations", "locations", elocs)

	elocs = gofn.Filter(elocs, func(loc string) bool {
		dir, err := utils.IsDirectory(loc)
		if err != nil || !dir {
			if os.IsNotExist(err) || !dir {
				log.Infof("Configuration location %s does not exist", loc)
			} else if !dir {
				log.Infof("Configuration location %s is not a directory", loc)
			} else {
				log.Infof("Configuration location %s cannot be searched", loc)
			}
			return false
		}
		return true
	})
	return elocs
}
