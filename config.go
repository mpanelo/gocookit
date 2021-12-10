package main

import "fmt"

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "",
		Name:     "gocookit_dev",
	}
}

func (c PostgresConfig) ConnectionInfo() string {
	if c.Password == "" {
		return fmt.Sprintf("host=%v user=%v dbname=%v port=%v sslmode=disable", c.Host, c.User, c.Name, c.Port)
	}

	return fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", c.Host, c.User, c.Password, c.Name, c.Port)
}

type Config struct {
	Port    int    `json:"port"`
	Env     string `json:"env"`
	Pepper  string `json:"pepper"`
	HMACKey string `json:"hmac_key"`
}

func DefaultConfig() Config {
	return Config{
		Port:    8000,
		Env:     "dev",
		Pepper:  "pepper",
		HMACKey: "secret",
	}
}

func (c Config) IsProd() bool {
	return c.Env == "prod"
}
