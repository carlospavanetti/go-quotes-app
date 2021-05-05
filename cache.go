package main

import (
	"log"
	"os"

	"github.com/go-redis/redis"
)

func CreateRedisClient() *redis.Client {
	var (
		host     = envVar("REDIS_HOST", "localhost")
		port     = envVar("REDIS_PORT", "6379")
		password = envVar("REDIS_PASSWORD", "")
	)

	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func envVar(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
