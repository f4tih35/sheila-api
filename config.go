package main

import "time"

type Config struct {
	RedisAddress  string
	RedisPassword string
	RedisDB       int
	TCPAddress    string
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
}

func LoadConfig() *Config {
	return &Config{
		RedisAddress:  "localhost:6379",
		RedisPassword: "",
		RedisDB:       0,
		TCPAddress:    ":8080",
		ReadTimeout:   5 * time.Minute,
		WriteTimeout:  5 * time.Minute,
	}
}
