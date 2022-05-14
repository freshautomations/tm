package tmconfig

import (
	"github.com/spf13/viper"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"tm/tm/v2/ux"
)

// Filename describes the tm config file path.
type Filename struct {
	Path      string
	Dir       string
	Base      string
	Extension string
	BaseNoExt string
}

// FindConfigFilename returns the filename properties of the tm configuration file.
func FindConfigFilename() Filename {
	var configFilePath string
	var err error
	var base string
	var ext string
	var dir string

	if viper.IsSet("config") {
		configFilePath, err = filepath.Abs(os.ExpandEnv(strings.TrimSpace(viper.GetString("config"))))
		if err != nil {
			ux.Fatal("could not find path to config %s", viper.GetString("config"))
		}
	} else {
		if viper.IsSet("home") {
			var configHome string
			configHome, err = filepath.Abs(os.ExpandEnv(strings.TrimSpace(viper.GetString("home"))))
			if err != nil {
				ux.Fatal("could not find path to config home %s", viper.GetString("home"))
			}
			configFilePath = filepath.Join(configHome, "config.toml")
		} else {
			configFilePath = os.ExpandEnv(filepath.FromSlash("$HOME/.tm/config.toml"))
			if err != nil {
				ux.Fatal("could not find default path to config %s", viper.GetString("config"))
			}
		}
	}
	dir = filepath.Dir(configFilePath)
	base = filepath.Base(configFilePath)
	ext = filepath.Ext(configFilePath)
	baseNoExt := strings.TrimSuffix(base, ext)
	if len(ext) > 0 {
		ext = ext[1:] // remove leading dot
	}
	return Filename{
		Path:      configFilePath,
		Dir:       dir,
		Base:      base,
		Extension: ext,
		BaseNoExt: baseNoExt,
	}
}

// CreateConfigPath searches for the configuration on the file system and creates a default one if it does not exist.
func CreateConfigPath() {
	// Find config
	filename := FindConfigFilename()
	ux.Debug("tm config %s", filename.Path)

	// Create config folder, if necessary.
	_, err := os.Stat(filename.Dir)
	if os.IsNotExist(err) {
		ux.Debug("creating tm config folder")
		err = os.MkdirAll(filename.Dir, fs.ModeDir|fs.ModePerm)
		if err != nil {
			ux.Fatal("could not create config file directory %s", filename.Dir)
		}
	} else if err != nil {
		ux.Fatal("could not create config file directory %s", filename.Dir)
	}
}
