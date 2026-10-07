package main

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	Port  string
	Store StoreConfig
}

type StoreConfig struct {
	DBURL     string
	CacheSize int
}

func LoadConfig() (Config, error) {
	cacheSize, err := strconv.Atoi(os.Getenv("CACHE_SIZE"))
	if err != nil {
		return Config{}, err
	}

	port := os.Getenv("PORT")
	if port == "" {
		return Config{}, errors.New("no value found for PORT environment variable")
	}
	dburl := os.Getenv("DATABASE_URL")
	if dburl == "" {
		return Config{}, errors.New("no value found for DATABASE_URL environment variable")
	}

	return Config{
		Port: port,
		Store: StoreConfig{
			DBURL:     dburl,
			CacheSize: cacheSize,
		},
	}, nil
}
