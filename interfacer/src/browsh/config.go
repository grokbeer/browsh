package browsh

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"bytes"

	"github.com/shibukawa/configdir"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	configFilename = "config.toml"

	isDebug = pflag.Bool("debug", false, "Log to ./debug.log")
	timeLimit = pflag.Int("time-limit", 0, "Kill Browsh after the specified number of seconds")
	_ = pflag.Bool("http-server-mode", false, "Run as an HTTP service")

	_ = pflag.String("startup-url", "https://google.com", "URL to launch at startup")
	_ = pflag.String("firefox.path", "firefox", "Path to Firefox executable")
	_ = pflag.Bool("firefox.with-gui", false, "Don't use headless Firefox")
	_ = pflag.Bool("firefox.use-existing", false, "Whether Browsh should launch Firefox or not")
)

func getConfigNamespace() string {
	if IsTesting {
		return "browsh-testing"
	}
	return "browsh"
}

// Gets a cross-platform path to a folder containing Browsh config
func getConfigDir() string {
	marker := "browsh-settings"
	// configdir has no other option but to have a nested folder
	configDirs := configdir.New(getConfigNamespace(), marker)
	folders := configDirs.QueryFolders(configdir.Global)
	// Delete the previously enforced nested folder
	path := strings.Trim(folders[0].Path, marker)
	os.MkdirAll(path, os.ModePerm)
	ensureConfigFile(path)
	return path
}

// Copy the sample config file if the user doesn't already have a config file
func ensureConfigFile(path string) {
	fullPath := filepath.Join(path, configFilename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		file, err := os.Create(fullPath)
		if err != nil {
			Shutdown(err)
		}
		defer file.Close()
		_, err = file.WriteString(configSample)
		if err != nil {
			Shutdown(err)
		}
	}
}

// Gets a cross-platform path to store a Browsh-specific Firefox profile
func getFirefoxProfilePath() string {
	configDirs := configdir.New(getConfigNamespace(), "firefox_profile")
	folders := configDirs.QueryFolders(configdir.Global)
	folders[0].MkdirAll()
	return folders[0].Path
}

func loadConfig() {
	dir := getConfigDir()
	fullPath := filepath.Join(dir, configFilename)
	Log("Looking in " + fullPath + " for config.")
	viper.SetConfigType("toml")
	viper.SetConfigName(strings.Trim(configFilename, ".toml"))
	viper.AddConfigPath(dir)
	viper.AddConfigPath(".")
	// First load the sample config in case the user hasn't updated any new fields
	err := viper.ReadConfig(bytes.NewBuffer([]byte(configSample)))
	if err != nil {
		Shutdown(err)
	}
	// Then load the users own config file, overwriting the sample config
	err = viper.MergeInConfig()
	if err != nil {
		Shutdown(err)
	}
	viper.BindPFlags(pflag.CommandLine)
	Log("Using the folowing config values:")
	Log(fmt.Sprintf("%v", viper.AllSettings()))
}