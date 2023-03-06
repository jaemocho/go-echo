package config

import (
	"flag"
	"path"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/creasty/defaults"
)

type Config struct {
	Listen        string `toml:"listen"`
	Phase         string `toml:"phase"`
	SqliteDBPath  string `toml:"sqliteDBPath"`
	JWTSigningKey string `toml:"jwtSigningKey"`
}

func New() (Config, error) {
	var configPath = ""

	_, currentFile, _, _ := runtime.Caller(0)
	currentDirectory := path.Dir(currentFile)

	flag.StringVar(&configPath, "config-file", currentDirectory+"/config.toml", "path to config file")
	flag.Parse()

	config := Config{}
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return Config{}, err
	}

	if err := defaults.Set(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}