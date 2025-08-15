package config

import (
	"flag"
	"os"
	"strconv"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

var path *Flags

func Load() {
	path = fetchConfigPath()
}

// Parse env file
// Return env data
func LoadEnv() *ConfigEnv {
	if err := godotenv.Load(path.Env); err != nil {
		panic(err)
	}
	redisNumDB, err := strconv.Atoi(os.Getenv("REDIS_NUM_DB"))
	if err != nil {
		panic(err)
	}
	return &ConfigEnv{
		NAME_DB:                    os.Getenv("NAME_DB"),
		USER_DB:                    os.Getenv("USER_DB"),
		PASS_DB:                    os.Getenv("PASS_DB"),
		PORT_DB:                    os.Getenv("PORT_DB"),
		HOST_DB:                    os.Getenv("HOST_DB"),
		HOST_REPLIC_DB:             os.Getenv("HOST_REPLIC_DB"),
		TYPE_DB:                    os.Getenv("TYPE_DB"),
		PRIVATE_KEY:                os.Getenv("PRIVATE_KEY"),
		PUBLIC_KEY:                 os.Getenv("PUBLIC_KEY"),
		RABBIT_NAME:                os.Getenv("RABBIT_NAME"),
		RABBIT_HOST:                os.Getenv("RABBIT_HOST"),
		RABBIT_PASSWORD:            os.Getenv("RABBIT_PASSWORD"),
		RABBIT_PORT:                os.Getenv("RABBIT_PORT"),
		RABBIT_QUEUE_FORGOT_TOKEN:  os.Getenv("RABBIT_QUEUE_FORGOT_TOKEN"),
		RABBIT_QUEUE_CONFIRM_EMAIL: os.Getenv("RABBIT_QUEUE_CONFIRM_EMAIL"),
		REDIS_HOST:                 os.Getenv("REDIS_HOST"),
		REDIS_PORT:                 os.Getenv("REDIS_PORT"),
		REDIS_PASSWORD:             os.Getenv("REDIS_PASSWORD"),
		REDIS_NUM_DB:               redisNumDB,
		HASH_PEPPER:                os.Getenv("HASH_PEPPER"),
	}
}

// Parse yaml file
// Return yaml data
func LoadYaml() *ConfigYaml {

	if path.Yaml == "" {
		panic("No read params `yaml`")
	}

	yaml := getYamlData(path.Yaml)

	return yaml
}

// Check flags
func fetchConfigPath() *Flags {

	var f Flags

	// read params config
	flag.StringVar(&f.Yaml, "yaml", "", "path to config file")
	flag.StringVar(&f.Env, "env", "", "path to data file")
	flag.Parse()

	if f.Yaml == "" {
		f.Yaml = os.Getenv("CONFIG_PATH")
		if f.Yaml == "" {
			f.Yaml = "./configs/config.yaml"
		}
	}

	if f.Env == "" {
		f.Env = os.Getenv("CONFIG_ENV")
		if f.Env == "" {
			f.Env = "./configs/config.env"
		}
	}

	return &f
}

func getYamlData(path string) *ConfigYaml {

	var res ConfigYaml

	if _, err := os.Stat(path); err != nil {
		panic("Does not exist file in dir: " + path)
	}
	if err := cleanenv.ReadConfig(path, &res); err != nil {
		panic("Error " + err.Error())
	}

	return &res
}
