package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Yaml struct {

}

type Config struct {
	Server 	ServerConfig
	DBMysql DBMysql
}

type ServerConfig struct {
	Host 	string
	Port 	int
}

type DBMysql struct {
	Host		string
	Port		string
	Username 	string
	Password 	string
	Database	string

}


var Conf Config

func (y *Yaml) LoadToml()  {

	yaml := viper.New()
	yaml.SetConfigName("app")
	yaml.SetConfigType("yaml")
	yaml.AddConfigPath("config")
	if err := yaml.ReadInConfig(); err != nil {
		fmt.Println(err)
		return
	}


	err := yaml.Unmarshal(&Conf)
	if err != nil {
		fmt.Println(err)
		return
	}




}