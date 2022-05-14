package config

import (
	"io/fs"
	"io/ioutil"
	"os"
	"tm/tm/v2/ux"
)

// Save tm config file to disk
func (cfg Config) Save() {
	bytes, err := cfg.CustomMarshal()
	if err != nil {
		ux.Fatal("could not encode config: %s", err)
	}
	// Write config
	err = ioutil.WriteFile(cfg.Filename.Path, bytes, fs.ModePerm)
	if err != nil {
		ux.Fatal("could not write config file %s: %s", cfg.Filename.Path, err.Error())
	}
}

// SaveNotOverwrite writes a new tm config file on disk, if it does not exist yet
func (cfg Config) SaveNotOverwrite() {
	_, err := os.Stat(cfg.Filename.Path)
	if os.IsNotExist(err) {
		ux.Debug("writing tm config file")
		cfg.Save()
	} else {
		ux.Debug("config file found")
	}
}
