package config

import "os"

// Connections
var RabbitMqUrl = os.Getenv("RABBITMQ_URL")
var DatabaseUrl = os.Getenv("DATABASE_URL")

// HTTP
var HttpPort = os.Getenv("HTTP_PORT")
var JwtSecret = os.Getenv("JWT_SECRET")

// Env
var Environment = os.Getenv("ENVIRONMENT")
var LogLevel = os.Getenv("LOG_LEVEL")

// Service
var ServiceName = "Follow Service"

// Queue
var UserEventsQueue = "follow.users.events"
