package config

import (
	"bufio"
	"encoding/json"
	"os"
)

type config struct {
	Env               string `json:"env"`
	Port			  string `json:"port"`
	Mysql             Mysql  `json:"mysql"`
	Redis             Redis  `json:"redis"`
	Log               Log    `json:"log"`
	ExcludeAuth 	map[string]map[string]bool		// 存放不校验的URL
}

type Auth struct{
	Name string
	Secret string
}

type Mysql struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
}

type Log struct {
	Level      string `json:"level"`
	Filename   string `json:"filename"`
	MaxSize    int    `json:"maxsize"`
	MaxAge     int    `json:"max_age"`
	MaxBackups int    `json:"max_backups"`
}

type Redis struct {
	Host string `json:"host"`
	Port string `json:"port"`
	DB   int    `json:"db"`
}

var (
	Config config
)

func init() {
	Config = config{}
	file, err := os.Open("conf/config.json")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(&Config); err != nil {
		panic(err)
	}

	Config.ExcludeAuth = map[string]map[string]bool{
		"GET": {"": true},
	}
}
