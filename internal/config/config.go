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
	App AppSettings `envPrefix:"app_"`
}

type AppSettings struct {
	Host string `env:"host"`
	Port int    `env:"port"`
}
