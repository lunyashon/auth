package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string        `yaml:"env" env-default:"local"`
	TokenJWT time.Duration `yaml:"token_jwt" env-default:"true"`
	GRPS     ConfigGRPC    `yaml:"grps"`
}

type ConfigGRPC struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {

	var res Config

	// exist config params
	path := fetchConfigPath()
	if path == "" {
		panic("No read params `config=path`")
	}

	// exist file path=./../../configs/local.yaml
	_, err := os.Stat(path)
	if err != nil {
		panic("Does not exist file in dir: " + path)
	}

	// writes *Config params
	if err := cleanenv.ReadConfig(path, &res); err != nil {
		panic("Error " + err.Error())
	}

	return &res
}

func fetchConfigPath() string {

	var res string

	// read params config
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	return res
}
