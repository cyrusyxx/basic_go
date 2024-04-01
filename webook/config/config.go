package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"path/filepath"
)

func InitConfig() {
	// Init Pflag
	configfile := pflag.StringP("config", "c",
		"config/config.yaml", "config file")
	pflag.Parse()

	// Init Viper
	viper.SetConfigFile(filepath.FromSlash(*configfile))
	fmt.Println("Config File:", *configfile)
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("Config file not found")
		}
		if _, ok := err.(viper.ConfigParseError); ok {
			panic("Config file parse error")
		}
		panic(err)
	}

	// Watch Config
	viper.OnConfigChange(func(e fsnotify.Event) {
		println("Config file changed:",
			e.Name, e.Op)
	})
	viper.WatchConfig()
}

// TODO maybe add a remote center of config
// TODO in 六.8
