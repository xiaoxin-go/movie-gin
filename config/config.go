package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Database Mysql `json:"database"`
	Server   struct {
		Port int `json:"port"`
	} `json:"server"`
}

type Mysql struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
}

var AppConfig Config

func LoadConfig(filePath string) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	err = json.Unmarshal(data, &AppConfig)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
}
