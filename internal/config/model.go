package config

import "time"

type ConfigYaml struct {
	Env             string        `yaml:"env" env-default:"local"`
	GRPS            ConfigGRPC    `yaml:"grps"`
	AccessTokenTTL  time.Duration `yaml:"token_access_ttl" env-default:"15m"`
	RefreshTokenTTL time.Duration `yaml:"token_refresh_ttl" env-default:"24h"`
	PathToLog       string        `yaml:"log_path"`
	NameSSOService  string        `yaml:"name_sso_service" env-default:"sso.domen.ru"`
	Rabbit          ConfigRabbit  `yaml:"rabbit"`
}

type ConfigRabbit struct {
	MaxRetries int           `yaml:"max_retries" env-default:"5"`
	RetryDelay time.Duration `yaml:"retry_delay" env-default:"1s"`
}

type ConfigGRPC struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type ConfigEnv struct {
	NAME_DB                    string `env:"NAME_DB" env-default:"sso"`
	USER_DB                    string `env:"USER_DB" env-default:"root"`
	PASS_DB                    string `env:"PASS_DB" env-default:"root"`
	PORT_DB                    string `env:"PORT_DB" env-default:"5432"`
	HOST_DB                    string `env:"HOST_DB" env-default:"localhost"`
	HOST_REPLIC_DB             string `env:"HOST_REPLIC_DB" env-default:"localhost"`
	TYPE_DB                    string `env:"TYPE_DB" env-default:"postgres"`
	PRIVATE_KEY                string `env:"PRIVATE_KEY" env-default:""`
	PUBLIC_KEY                 string `env:"PUBLIC_KEY" env-default:""`
	RABBIT_NAME                string `env:"RABBIT_NAME" env-default:"guest"`
	RABBIT_HOST                string `env:"RABBIT_HOST" env-default:"localhost"`
	RABBIT_PASSWORD            string `env:"RABBIT_PASSWORD" env-default:"guest"`
	RABBIT_PORT                string `env:"RABBIT_PORT" env-default:"5672"`
	RABBIT_QUEUE_FORGOT_TOKEN  string `env:"RABBIT_QUEUE_FORGOT_TOKEN" env-default:"forgot_token"`
	RABBIT_QUEUE_CONFIRM_EMAIL string `env:"RABBIT_QUEUE_CONFIRM_EMAIL" env-default:"confirm_email"`
	REDIS_HOST                 string `env:"REDIS_HOST" env-default:"localhost"`
	REDIS_PORT                 string `env:"REDIS_PORT" env-default:"6379"`
	REDIS_PASSWORD             string `env:"REDIS_PASSWORD" env-default:""`
	REDIS_NUM_DB               int    `env:"REDIS_NUM_DB" env-default:"1"`
	HASH_PEPPER                string `env:"HASH_PEPPER" env-default:"e-toolnet"`
}

type Flags struct {
	Yaml string
	Env  string
}

type Config struct {
	ConfigEnv  ConfigEnv
	ConfigYaml ConfigYaml
}
