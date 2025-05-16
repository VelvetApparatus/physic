package config

import "github.com/caarlos0/env/v7"

var c *Config

func C() *Config { return c }

func Init() error {
	c = new(Config)
	err := env.Parse(c)
	if err != nil {
		return err
	}
	return nil
}

type Config struct {
	App AppSettings   `envPrefix:"app_"`
	Db  DBSettings    `envPrefix:"DB_"`
	S3  MinioSettings `envPrefix:"S3_"`
}

type AppSettings struct {
	Host string `env:"host"`
	Port int    `env:"port"`
}

type DBSettings struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Name     string `env:"NAME"`
	SSLMode  string `env:"SSL_MODE"`
	AppName  string `env:"APP_NAME"`
}

type MinioSettings struct {
	Host            string `env:"HOST"`
	Port            string `env:"PORT"`
	AccessKeyID     string `env:"ACCESS_KEY_ID"`
	SecretAccessKey string `env:"SECRET_ACCESS_KEY"`
	UseSSL          bool   `env:"USE_SSL"`
}
