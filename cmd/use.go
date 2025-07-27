package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/boazj/muxocil/common"
	"github.com/boazj/muxocil/provider"
	"github.com/boazj/muxocil/utils"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use",
	Short: "Use a layout, specified by name or path",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := common.GetConfig()
		if len(args) == 0 {
			log.Fatal("Parameter missing to execute command - Use command missing actual layout")
			utils.ExitError(utils.ExitCmdUseBadCommand)
		}
		layoutPath := getLayoutPath(cfg, args[0])
		opts := &common.CommandOpts{
			Here:  false,
			Show:  false,
			Debug: false,
		}

		// TODO: from env and getprovider needs to verify if you can run a multiplex program if launch_multiplexer_app is on
		prov := getProvider(cfg, opts)
		err := provider.Process(prov, layoutPath)
		if err != nil {
			// TODO: log and choose a better error code
			utils.ExitError(utils.ExitProviderFailure)
		}
	},
}

func getProvider(cfg *common.Config, opts *common.CommandOpts) common.Provider {
	if cfg.UseProviderFailFast || (!cfg.UseProviderFailFast && !cfg.InfferProviderFromEnv) {
		log.Debug(
			"Fail fast mode",
			"config failfast",
			cfg.UseProviderFailFast,
			"not failfast not inffer",
			!cfg.UseProviderFailFast && !cfg.InfferProviderFromEnv,
		)
		mux, ok := provider.ProviderDefs.GetProvider(common.MuxID(cfg.UseProvider))
		if !ok {
			log.Error("Unsupported use provider configuration", "provider", cfg.UseProvider)
			utils.ExitError(utils.ExitProviderFailFast)
		}
		prov, err := provider.NewProvider(cfg, mux.ID, opts)
		if err != nil {
			log.Error("Failed to create provider", "provider", cfg.UseProvider, "error", err)
			utils.ExitError(utils.ExitProviderFailFast)
		}

		if cfg.ValidateProviderFromEnv {
			envProv, err := provider.FromEnv(cfg, opts)
			if err != nil {
				log.Error("Failed to create provider", "error", err)
				utils.ExitError(utils.ExitProviderFailFast)
			}
			if envProv.GetID() != prov.GetID() {
				log.Error(
					"Provider validation failed, mismatch between configured provider and the environment",
					"provider",
					prov.GetID(),
					"environment",
					envProv.GetID(),
				)
				utils.ExitError(utils.ExitProviderValidation)
			}
			log.Info("Provider validation successful", "provider", cfg.UseProvider)
		}
		log.Info("Provider creation successful", "provider", cfg.UseProvider)
		return prov
	}
	mux, ok := provider.ProviderDefs.GetProvider(common.MuxID(cfg.UseProvider))
	if !ok {
		log.Info("Unsupported use provider configuration, falling back to inffered provider", "provider", cfg.UseProvider)
		envProv, err := provider.FromEnv(cfg, opts)
		if err != nil {
			log.Error("Failed to create fallback provider", "error", err)
			utils.ExitError(utils.ExitProviderFailure)
		}
		return envProv
	}
	prov, err := provider.NewProvider(cfg, mux.ID, opts)
	if err != nil {
		log.Error("Failed to create provider, falling back to inffered provider", "error", err)
		envProv, err := provider.FromEnv(cfg, opts)
		if err != nil {
			log.Error("Failed to create fallback provider", "error", err)
			utils.ExitError(utils.ExitProviderFailure)
		}
		log.Info("Fallback provider creation successful", "provider", envProv.GetID())
		return envProv
	}
	log.Info("Provider creation successful", "provider", cfg.UseProvider)
	return prov
}

func getLayoutPath(cfg *common.Config, candidate string) string {
	isName := !utils.IsYaml(candidate)
	isFilename := utils.IsYaml(candidate) && !strings.Contains(candidate, string(os.PathSeparator))
	isFilepath := utils.IsYaml(candidate) && strings.Contains(candidate, string(os.PathSeparator))

	if !isName && !isFilename && !isFilepath {
		log.Fatal("Invalid parameter cannot be handled")
		utils.ExitError(utils.ExitCmdUseBadCommand)
	}
	layout := ""
	if isFilename {
		layout = os.ExpandEnv(candidate)
	} else {
		for _, loc := range cfg.GetLayoutSearchLocations() {
			layouts := utils.GetFilesRecursively(loc, utils.IsYaml)
			for _, l := range layouts {
				if candidate == filepath.Base(l) || candidate == utils.GetFileNamePart(l) {
					layout = l
					break
				}
			}
			if layout != "" {
				break
			}
		}
	}
	if layout == "" {
		log.Fatal("Provided layout cannot be recognized")
		utils.ExitError(utils.ExitCmdUseBadCommand)
	}
	if _, err := os.Stat(layout); err != nil {
		log.Fatal("Cannot use layout file", "file", layout, "error", err)
		utils.ExitError(utils.ExitCmdUseBadFile)
	}
	return layout
}

func init() {
	rootCmd.AddCommand(useCmd)
}
