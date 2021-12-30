package config

import (
	"encoding/json"
	"os"

	"go.uber.org/zap"
)

type MongoConfiguration struct {
	URI        string `json:"uri"`
	Database   string `json:"database"`
	Collection string `json:"collection"`
}

type Configuration struct {
	Mongo          MongoConfiguration `json:"mongo"`
	ConsumerKey    string             `json:"consumerKey"`
	ConsumerSecret string             `json:"consumerSecret"`
	AccessToken    string             `json:"accessToken"`
	AccessSecret   string             `json:"accessSecret"`
}

func ReadConfig(logger *zap.Logger) Configuration {
	file, err := os.Open("config.json")
	if err != nil {
		logger.Fatal("Can not open configuration file.", zap.String("originalError", err.Error()))
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var config Configuration
	err = decoder.Decode(&config)
	if err != nil {
		logger.Fatal("Can not decode configuration file.", zap.String("originalError", err.Error()))
	}
	return config
}
