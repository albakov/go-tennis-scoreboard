package config

import (
	"github.com/BurntSushi/toml"
	"os"
	"path"
)

type Config struct {
	Host     string `toml:"host"`
	Port     int64  `toml:"port"`
	PathToDB string `toml:"abs_path_to_database"`
	Game
}

type Game struct {
	MaxSetsToWin   int64 `toml:"max_sets_to_win"`
	MatchesPerPage int64 `toml:"matches_per_page"`
}

func MustNew() *Config {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	c := &Config{}

	_, err = toml.DecodeFile(path.Join(wd, "config", "app.toml"), c)
	if err != nil {
		panic(err)
	}

	return c
}
