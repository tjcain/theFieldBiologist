package main

import "fmt"

// PostgresConfig holds config variables for connecting to postgres
type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// DefaultPostgresConfig generates a PostgresConfig with default values
func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "",
		Name:     "fieldbiologist_dev",
	}
}

// Dialect returns the gorm dialect in use
func (c PostgresConfig) Dialect() string {
	return "postgres"
}

// ConnectionInfo returns a database connection string
func (c PostgresConfig) connectionInfo() string {
	if c.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname =%s "+
			"sslmode=disable", c.Host, c.Port, c.User, c.Name)
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname =%s "+
		"sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Name)
}

// Config holds start up configuration variables
type Config struct {
	Port    int    `json:"port"`
	Env     string `json:"env"`
	Pepper  string `json:"pepper"`
	HMACKey string `json:"hmac_key"`
}

// DefaultConfig generates development environment variables
func DefaultConfig() Config {
	return Config{
		Port:    8080,
		Env:     "dev",
		Pepper:  "secret-random-string",
		HMACKey: "secret-hmac-key",
	}
}

// IsProd returns true if app is build for production
func (c Config) IsProd() bool {
	return c.Env == "prod"
}
