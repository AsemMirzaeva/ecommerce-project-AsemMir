package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
    "github.com/spf13/viper"
)

type PostgresConfig  struct{
	Host string 
	Port string 
	User string 
	Password string 
	Database string
}


type Config struct {
    Postgres PostgresConfig

    APIGatewayHost     string
    APIGatewayPort     string
    UserServiceHost    string
    UserServicePort    string
    ProductServiceHost string
    ProductServicePort string
    OrderServiceHost   string
    OrderServicePort   string
}

func LoadConfig(path string) (*Config, error) {
    err := godotenv.Load(path + "/.env")
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    cfg := viper.New()
	cfg.AutomaticEnv()

    config := &Config{
		Postgres: PostgresConfig{
			
			Host: cfg.GetString("POSTGRES_HOST"),
			Port: cfg.GetString("POSTGRES_PORT"),
			User:  cfg.GetString("POSTGRES_USER"),
			Password: cfg.GetString("POSTGRES_PASSWORD"),
			Database:  cfg.GetString("POSTGRES_DATABASE"),
		},
		APIGatewayHost:     os.Getenv("APIGATEWAY_HOST"),
        APIGatewayPort:     os.Getenv("APIGATEWAY_PORT"),

        UserServiceHost:    os.Getenv("USERSERVICE_HOST"),
        UserServicePort:    os.Getenv("USERSERVICE_PORT"),

        ProductServiceHost: os.Getenv("PRODUCTSERVICE_HOST"),
        ProductServicePort: os.Getenv("PRODUCTSERVICE_PORT"),

        OrderServiceHost:   os.Getenv("ORDERSERVICE_HOST"),
        OrderServicePort:   os.Getenv("ORDERSERVICE_PORT"),

	}

    return config, nil
}
