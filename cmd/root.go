package cmd

import (
	"fmt"

	"github.com/autobrr/tqm/pkg/config"
	"github.com/autobrr/tqm/pkg/logger"
	"github.com/autobrr/tqm/pkg/runtime"
	"github.com/autobrr/tqm/pkg/stringutils"
	"github.com/autobrr/tqm/pkg/tracker"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	// Global flags

	FlagLogLevel     = 0
	FlagConfigFile   = "config.yaml"
	FlagConfigFolder = config.GetDefaultConfigDirectory("tqm", FlagConfigFile)
	FlagLogFile      = "activity.log"

	//flagFilterName                       string

	FlagDryRun                           bool
	FlagExperimentalRelabelForCrossSeeds bool

	// Global vars
	Log         *logrus.Entry
	Initialized bool
)

//var rootCmd = &cobra.Command{
//	Use:   "tqm",
//	Short: "A CLI torrent queue manager",
//	Long: `A CLI application that can be used to manage your torrent clients.
//`,
//}
//
//func Execute() {
//	if err := rootCmd.Execute(); err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//}

//func init() {
//	// Parse persistent flags
//	rootCmd.PersistentFlags().StringVar(&FlagConfigFolder, "config-dir", FlagConfigFolder, "Config folder")
//	rootCmd.PersistentFlags().StringVarP(&FlagConfigFile, "config", "c", FlagConfigFile, "Config file")
//	rootCmd.PersistentFlags().StringVarP(&FlagLogFile, "log", "l", FlagLogFile, "Log file")
//	rootCmd.PersistentFlags().CountVarP(&FlagLogLevel, "verbose", "v", "Verbose level")
//
//	rootCmd.PersistentFlags().BoolVar(&FlagDryRun, "dry-run", false, "Dry run mode")
//	rootCmd.PersistentFlags().BoolVar(&FlagExperimentalRelabelForCrossSeeds, "experimental-relabel", false, "Enable experimental relabeling for cross-seeded torrents, using hardlinks (only qbit for now")
//}

func initCore(showAppInfo bool) {
	// Set core variables
	//if !rootCmd.PersistentFlags().Changed("config") {
	//	FlagConfigFile = filepath.Join(FlagConfigFolder, FlagConfigFile)
	//}
	//if !rootCmd.PersistentFlags().Changed("log") {
	//	FlagLogFile = filepath.Join(FlagConfigFolder, FlagLogFile)
	//}

	// Init Logging
	if err := logger.Init(FlagLogLevel, FlagLogFile); err != nil {
		Log.WithError(err).Fatal("Failed to initialize logging")
	}

	Log = logger.GetLogger("app")

	// Init Config
	if err := config.Init(FlagConfigFile); err != nil {
		Log.WithError(err).Fatal("Failed to initialize config")
	}

	// Init Trackers
	if err := tracker.Init(config.Config.Trackers); err != nil {
		Log.WithError(err).Fatal("Failed to initialize trackers")
	}

	// Show App Info
	if showAppInfo {
		showUsing()
	}
}

func showUsing() {
	// show app info
	Log.Infof("Using %s = %s (%s@%s)", stringutils.LeftJust("VERSION", " ", 10),
		runtime.Version, runtime.GitCommit, runtime.Timestamp)
	logger.ShowUsing()
	config.ShowUsing()
	Log.Info("------------------")
}

func validateClientEnabled(clientConfig map[string]interface{}) error {
	v, ok := clientConfig["enabled"]
	if !ok {
		return fmt.Errorf("no enabled setting found in client configuration: %+v", clientConfig)
	} else {
		enabled, ok := v.(bool)
		if !ok || !enabled {
			return errors.New("client is not enabled")
		}
	}

	return nil
}

func getClientConfigString(setting string, clientConfig map[string]interface{}) (*string, error) {
	v, ok := clientConfig[setting]
	if !ok {
		return nil, fmt.Errorf("no %q setting found in client configuration: %+v", setting, clientConfig)
	}

	value, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("failed type-asserting %q of client: %#v", setting, v)
	}

	return &value, nil
}

func getClientDownloadPathMapping(clientConfig map[string]interface{}) (map[string]string, error) {
	v, ok := clientConfig["download_path_mapping"]
	if !ok {
		return nil, nil
	}

	tmp, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed type-asserting download_path_mapping of client: %#v", v)
	}

	clientDownloadPathMapping := make(map[string]string)
	for k, v := range tmp {
		if vv, ok := v.(string); ok {
			clientDownloadPathMapping[k] = vv
		} else {
			return nil, fmt.Errorf("failed type-asserting download_path_mapping of client for %q: %#v", k, v)
		}
	}

	return clientDownloadPathMapping, nil
}

func getClientFilter(clientConfig map[string]interface{}) (*config.FilterConfiguration, error) {
	v, ok := clientConfig["filter"]
	if !ok {
		return nil, fmt.Errorf("no filter setting found in client configuration: %+v", clientConfig)
	}

	clientFilterName, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("failed type-asserting filter of client: %#v", v)
	}

	clientFilter, ok := config.Config.Filters[clientFilterName]
	if !ok {
		return nil, fmt.Errorf("failed finding configuration of filter: %+v", clientFilterName)
	}

	return &clientFilter, nil
}

func getFilter(filterName string) (*config.FilterConfiguration, error) {
	clientFilter, ok := config.Config.Filters[filterName]
	if !ok {
		return nil, fmt.Errorf("failed finding configuration of filter: %+v", filterName)
	}

	return &clientFilter, nil
}
