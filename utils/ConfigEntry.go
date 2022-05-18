package utils

// ConfigEntry gets and sets TOML and JSON configuration entries for

import (
	"github.com/spf13/viper"
	"strings"
	"tm/tm/v2/ux"
)

func SetConfigEntry(file string, key string, value interface{}) {
	toml := viper.New()
	toml.SetConfigFile(file)
	err := toml.ReadInConfig()
	if err != nil {
		ux.Fatal("%s, cannot set key %s", err.Error(), key)
	}
	toml.Set(key, value)
	err = toml.WriteConfig()
	if err != nil {
		ux.Fatal("%s, cannot set key %s", err.Error(), key)
	}
}

func GetConfigEntry(file string, key string) interface{} {
	toml := viper.New()
	toml.SetConfigFile(file)
	err := toml.ReadInConfig()
	if err != nil {
		ux.Fatal("%s, cannot get key %s", err.Error(), key)
	}
	return toml.Get(key)
}

func GetConfigEntryContentString(content string, contentType string, key string) string {
	toml := viper.New()
	toml.SetConfigType(contentType)
	err := toml.ReadConfig(strings.NewReader(content))
	if err != nil {
		ux.Fatal("%s, cannot get key %s", err.Error(), key)
	}
	return toml.GetString(key)
}
