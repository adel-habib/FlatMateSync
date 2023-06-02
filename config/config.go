package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		Host           string
		Port           string
		User           string
		Password       string
		Name           string
		Migrations_Url string
	}

	Server struct {
		Host string
		Port string
	}

	Oauth struct {
		Client_Id     string
		Client_Secret string
	}
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
