package config

var c *Config

func C() *Config { return c }

func init() {
	c = new(Config)

}

type Config struct {
	App AppSettings `envPrefix:"app_"`
}

type AppSettings struct {
	Host string `env:"host"`
	Port int    `env:"port"`
}
