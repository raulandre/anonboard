package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	config *viper.Viper
}

func NewConfig() *Config {
	c := new(Config)
	c.config = readConfig()
	return c
}

func (c *Config) Get() *viper.Viper {
	if c.config == nil {
		logrus.Fatal("attempted to get non initialized config")
	}
	return c.config
}

func readConfig() *viper.Viper {
	logrus.Info("reading env")
	v := viper.New()
	v.AutomaticEnv()

	env := v.GetString("ENVIROMENT")
	if env == "" {
		env = "local"
	}

	logrus.Infof("env: %s", env)
	v.SetConfigName(env)
	v.AddConfigPath("config")
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		logrus.Fatalf("error reading env: %s", err.Error())
	}

	return v
}
