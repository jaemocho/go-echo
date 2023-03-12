package config

import (
	"flag"
	"path"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/creasty/defaults"
)

type Postgre struct {
	IP       string `toml:"ip"`
	Port     string `toml:"port"`
	DBName   string `toml:"dbname"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	SSLMode  string `toml:"sslmode"`
	TimeZone string `toml:"timeZone"`
}

type Config struct {
	Listen        string `toml:"listen"`
	Phase         string `toml:"phase"`
	JWTSigningKey string `toml:"jwtSigningKey"`

	GitClient   string `toml:"gitClient"`
	GitHubToken string `toml:"githubToken"`
	GitLabToken string `toml:"gitlabToken"`

	DB           string `toml:"db"`
	SqliteDBPath string `toml:"sqliteDBPath"`
	Postgre      `toml:"postgre"`
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
