package config

import (
	"fmt"
	"log"
	"reflect"

	"github.com/spf13/viper"
)

type Config struct {
	Environment     string `mapstructure:"ENVIRONMENT"`
	Host            string `mapstructure:"HOST"`
	AppVersion      string `mapstructure:"APP_VERSION"`
	QueueSystem     string `mapstructure:"QUEUE_SYSTEM"`
	MongoURI        string `mapstructure:"MONGO_URI"`
	MongoDB         string `mapstructure:"MONGO_DB"`
	MongoAuthSource string `mapstructure:"MONGO_AUTH_SOURCE"`
	MongoUser       string `mapstructure:"MONGO_USER"`
	MongoPassword   string `mapstructure:"MONGO_PASSWORD"`
	RedisHost       string `mapstructure:"REDIS_HOST"`
	RedisPort       string `mapstructure:"REDIS_PORT"`
	RedisPass       string `mapstructure:"REDIS_PASSWORD"`
	ServerPort      string `mapstructure:"APP_PORT"`
}

var AppConfig Config

// LoadConfig loads the config from environment variables
func InitializeConfig() {

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Error reading config file: %s", err)
		log.Println("Using environment variables only")
	}

	// Set default values if not provided
	viper.SetDefault("PORT", 8080)
	viper.SetDefault("HOST", "0.0.0.0")
	viper.SetDefault("ENVIRONMENT", "dev")
	viper.SetDefault("APP_VERSION", "0.0.1")

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		panic(fmt.Errorf("fatal error when unmarshaling config: %s", err))
	}

	t := reflect.TypeOf(AppConfig)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envKey := field.Tag.Get("mapstructure")
		if envKey != "" && envKey != "PORT" {
			envVal := viper.GetString(envKey)
			if envVal != "" {
				viper.Set(envKey, envVal)
			}
		}
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		panic(fmt.Errorf("fatal error when re-unmarshaling config: %s", err))
	}

	log.Printf("Configuration loaded. Environment: %s, Host: %s, Port: %d", AppConfig.Environment, AppConfig.Host, AppConfig.ServerPort)
}
