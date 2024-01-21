package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string        `mapstructure:"PORT"`
	ConnString  string        `mapstructure:"CONN_STRING"`
	DriverName  string        `mapstructure:"DRIVER_NAME"`
	Timeout     time.Duration `mapstructure:"TIMEOUT"`
	IdleTimeout time.Duration `mapstructure:"IDLE_TIMEOUT"`
}

func MustConfig() *Config {
	var config Config
	viper.AddConfigPath("./internal/config/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic("failed to read config: " + err.Error())
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err.Error())
	}
	return &config
}
