package config

import "os"

type Config struct {
	DatabaseUrl        string
	RabbitMqUrl        string
	SocmedExchangeName string
	DeadLetterExchangeName string
	UserCreatedEvent   string
	UserUpdatedEvent   string
	AllUserEvents   string
	UserEventsQueue string
		DeadLetterQueue string
}

func NewSettings() *Config {
	return &Config{
		DatabaseUrl:        os.Getenv("DATABASE_URL"),
		RabbitMqUrl:        os.Getenv("RABBITMQ_URL"),
		SocmedExchangeName: "socmed.events",
		DeadLetterExchangeName: "socmed.events.dlx",
		UserCreatedEvent:   "user.created",
		UserUpdatedEvent:   "user.updated",
		AllUserEvents: "user.*",
		UserEventsQueue:   "follow.user.queue",
		DeadLetterQueue: "socmed.events.dlq",
	}
}

var Settings = NewSettings()
