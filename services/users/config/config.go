package config

import "os"

type Config struct {
	DatabaseUrl        string
	RabbitMqUrl        string
	SocmedExchangeName string
	UserCreatedEvent   string
	UserUpdatedEvent   string
}

func NewSettings() *Config {
	return &Config{
		DatabaseUrl:        os.Getenv("DATABASE_URL"),
		RabbitMqUrl:        os.Getenv("RABBITMQ_URL"),
		SocmedExchangeName: "socmed.events",
		UserCreatedEvent:   "user.created",
		UserUpdatedEvent:   "user.updated",
	}
}

var RabbitMqUrl = os.Getenv("RABBITMQ_URL")
var DatabaseUrl = os.Getenv("DATABASE_URL")
var ServiceName = "Users Service"

var Settings = NewSettings()
