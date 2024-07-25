package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env string `yaml:"env" env-default:"local"`
}

func mustLoadPath(configPath string) *Config {
	if _, err := os.Stat(configPath); err != nil {
		e, _ := os.Executable()
		panic("Config does not exist on: " + e + configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("Something wrong with config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG")
	}

	return res
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("Config path is empty!")
	}

	return mustLoadPath(configPath)
}
