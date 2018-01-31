package config

import (
	"flag"
	"go-lambda-graphql/services/generate"
)

// Port defines server listening port
var Port string

// IsProduction describes server development mode
var IsProduction bool

// ConnectionString from db connection string
var ConnectionString string

// JWT secret
var JWTSecret string

// Directory represents http fileserver directory
var Directory string

func init() {
	Directory = *flag.String("d", "webapp/build", "the directory of static file to host")
	JWTSecret = generate.GenerateRandomString(64)
	ConnectionString = "user=williamhuang dbname=lambda sslmode=disable"
	IsProduction = *flag.Bool("p", false, "production mode?")
	Port = *flag.String("port", "3001", "listening port")
}
