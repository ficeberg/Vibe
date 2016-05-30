package models

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"time"
)

type Config struct {
	Title   string
	Build   string
	Owner   ownerInfo
	DB      database `toml:"database"`
	Servers map[string]server
	JWT     jwt
	List    list
	Social  map[string]social
}

type ownerInfo struct {
	Name string
	Org  string `toml:"organization"`
	Bio  string
	DOB  time.Time
}

type database struct {
	Host    string
	Ports   []int
	Name    string
	ConnMax int `toml:"connection_max"`
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
	conf := Config{}
	if _, err := toml.DecodeFile("./config.toml", &conf); err != nil {
		fmt.Println(err)
	}
	return conf
}
