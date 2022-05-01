package config

import (
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
	"tm/m/v2/ux"
)

// Filename describes the config file path.
type Filename struct {
	Path      string
	Dir       string
	Base      string
	Extension string
	BaseNoExt string
}

// FindConfigFilename returns the filename properties of the manager configuration file
func FindConfigFilename() Filename {
	var err error

	configFilePath := viper.GetString("config")
	if viper.IsSet("config") {
		configFilePath, err = filepath.Abs(configFilePath)
		if err != nil {
			ux.Fatal("could not find path %s", viper.GetString("config"))
		}
	} else {
		configHome := viper.GetString("home")
		configHome, err = filepath.Abs(configHome)
		if err != nil {
			ux.Fatal("could not find path %s", viper.GetString("home"))
		}
		configFilePath = filepath.Join(configHome, "config.toml")
	}
	base := filepath.Base(configFilePath)
	ext := filepath.Ext(configFilePath)
	baseNoExt := strings.TrimSuffix(base, ext)
	if len(ext) > 0 {
		ext = ext[1:] // remove leading dot
	}
	return Filename{
		Path:      configFilePath,
		Dir:       filepath.Dir(configFilePath),
		Base:      base,
		Extension: ext,
		BaseNoExt: baseNoExt,
	}
}
