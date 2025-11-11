package config

import "os"

type Config struct {
	DatabaseUrl        string
	RabbitMqUrl        string
	SocmedExchangeName string
	UserCreatedEvent   string
	UserUpdatedEvent   string
	AllUserEvents   string
	UserEventsQueue string
}

func NewSettings() *Config {
	return &Config{
		DatabaseUrl:        os.Getenv("DATABASE_URL"),
		RabbitMqUrl:        os.Getenv("RABBITMQ_URL"),
		SocmedExchangeName: "socmed.events",
		UserCreatedEvent:   "user.created",
		UserUpdatedEvent:   "user.updated",
		AllUserEvents: "user.*",
		UserEventsQueue:   "follow.user.queue",
	}
}

var Settings = NewSettings()
