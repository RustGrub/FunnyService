package config

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

type Config struct {
	AppName     string `toml:"AppName"`
	AppPort     string `toml:"AppPort"`
	Environment string `toml:"Environment"`
	ProjectsDb  *DatabaseCfg
	Redis       *RedisCfg
	Logger      *LoggerCfg
	Nats        *NatsCfg
	House       *DatabaseCfg
}

type RedisCfg struct {
	Server   string `toml:"Server"`
	Port     string `toml:"Port"`
	Database int    `toml:"Database"`
}

type DatabaseCfg struct {
	Server   string `toml:"Server"`
	Port     string `toml:"Port"`
	Database string `toml:"Database"`
	Username string
	Password string
}
type NatsCfg struct {
	Server string `toml:"Server"`
	Port   string `toml:"Port"`
}

type LoggerCfg struct {
	Level string `toml:"Level"`
}

func LoadConfig() *Config {
	conf := &Config{}

	configFile := "config/conf_local.toml"

	if _, err := toml.DecodeFile(configFile, conf); err != nil {
		log.Fatal("couldn't decode config file:", err)
	}
	conf.ProjectsDb.Username = os.Getenv("pg_user")
	conf.ProjectsDb.Password = os.Getenv("pg_password")
	conf.House.Username = "user"
	conf.House.Password = "password"
	/*
		conf.House.Username = os.Getenv("ch_user")
		conf.House.Password = os.Getenv("ch_password")
	*/

	// + redis pass...
	return conf
}
