package config

import (
	"fmt"
	"log"
	"reflect"

	"github.com/spf13/viper"
)

type Config struct {
	Environment   string `mapstructure:"ENVIRONMENT"`
	Host          string `mapstructure:"HOST"`
	AppVersion    string `mapstructure:"APP_VERSION"`
	QueueSystem   string `mapstructure:"QUEUE_SYSTEM"`
	MongoURI      string `mapstructure:"MONGO_URI"`
	MongoDB       string `mapstructure:"MONGO_DB"`
	MongoUser     string `mapstructure:"MONGO_USER"`
	MongoPassword string `mapstructure:"MONGO_PASSWORD"`
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPass     string `mapstructure:"REDIS_PASSWORD"`
	ServerPort    string `mapstructure:"APP_PORT"`
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

	// Override config with environment variables
	t := reflect.TypeOf(AppConfig)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envKey := field.Tag.Get("mapstructure")
		if envKey != "" && envKey != "PORT" { // Skip PORT as we've already handled it
			envVal := viper.GetString(envKey)
			if envVal != "" {
				viper.Set(envKey, envVal)
			}
		}
	}

	// Re-unmarshal to ensure all values are updated
	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		panic(fmt.Errorf("fatal error when re-unmarshaling config: %s", err))
	}

	log.Printf("Configuration loaded. Environment: %s, Host: %s, Port: %d", AppConfig.Environment, AppConfig.Host, AppConfig.ServerPort)
}

// // InitRedis initializes a Redis client
// func InitRedis(cfg Config) *redis.Client {
// 	redisAddr := cfg.RedisAddr
// 	redisPassword := cfg.RedisPass
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     redisAddr,
// 		Password: redisPassword, // No password set
// 		DB:       0,             // Use default DB
// 	})

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	if err := rdb.Ping(ctx).Err(); err != nil {
// 		log.Fatalf("Failed to connect to Redis: %v", err)
// 	}

// 	log.Println("Connected to Redis!")
// 	return rdb
// }

// func InitRabbitMQ() *amqp.Connection {
// 	rabbitmqURI := os.Getenv("RABBITMQ_URI")
// 	conn, err := amqp.Dial(rabbitmqURI)
// 	if err != nil {
// 		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
// 	}

// 	log.Println("Connected to RabbitMQ")
// 	return conn
// }
