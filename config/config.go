package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DB    DBConfig
	Kafka KafkaConfig
}

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

type KafkaConfig struct {
	URL   string
	Topic string
}

var Cfg Config

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&Cfg)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}
