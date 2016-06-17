package models

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Title   string
	Build   string
	Owner   ownerInfo
	DB      database `mapstructure:"database"`
	Servers map[string]server
	JWT     jwt
	List    list
	Social  map[string]social
}

type ownerInfo struct {
	Name string
	Org  string `mapstructure:"organization"`
	Bio  string
	DOB  time.Time
}

type database struct {
	Host    string
	Ports   []int
	Name    string
	ConnMax int `mapstructure:"connection_max"`
	Enabled bool
	Table   map[string]string
}

type server struct {
	IP   string
	DC   string
	Port string
}

type jwt struct {
	SigningKey    string
	SigningMethod string
	Bearer        string
	TokenTTL      time.Duration
}

type list struct {
	White []string
	Black []string
}

type social struct {
	Key    string
	Secret string
}

func (c Config) Init() Config {
	var vp = viper.New()

	vp.SetConfigName("config")
	vp.AddConfigPath(".")
	vp.AddConfigPath("$GOPATH/")
	vp.SetConfigType("toml")
	err := vp.ReadInConfig()
	if err != nil {
		fmt.Errorf("Fatal error config file: %s \n", err)
	}
	vp.Unmarshal(&c)

	return c
}
