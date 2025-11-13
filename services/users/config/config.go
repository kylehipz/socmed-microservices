package config

import "os"

var RabbitMqUrl = os.Getenv("RABBITMQ_URL")
var DatabaseUrl = os.Getenv("DATABASE_URL")
var HttpPort = os.Getenv("HTTP_PORT")
var JwtSecret = os.Getenv("JWT_SECRET")
var ServiceName = "Users Service"
var SocmedExchangeName = "socmed.events"

