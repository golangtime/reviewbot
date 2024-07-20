package config

import (
	"io"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Email     EmailConfig     `toml:"email"`
	Pachca    PachcaConfig    `toml:"pachca"`
	Bitbucket BitbucketConfig `toml:"bitbucket"`
	Schedule  ScheduleConfig  `toml:"schedule"`
}

type BitbucketConfig struct {
	URL      string `toml:"url"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}

type EmailConfig struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	From     string `toml:"from"`
}

type PachcaConfig struct {
	Token string `toml:"token"`
}

type ScheduleConfig struct {
	TestMode bool     `toml:"test_mode"`
	Default  [][3]int `toml:"default"`
}

func Load() (*Config, error) {
	f, err := os.Open("config.toml")
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var config *Config

	if err := toml.Unmarshal(body, &config); err != nil {
		return nil, err
	}

	return config, err
}
