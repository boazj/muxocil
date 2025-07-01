package common

import "github.com/spf13/viper"

const (
	ConfSearchLocations               = "configuration.layout_search_locations"
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

	viper.SetDefault(ConfSearchLocations, []string{"$Home/.teamocil"})
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
