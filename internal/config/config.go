package config

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
	cacheSizeRaw := os.Getenv("CACHE_SIZE")
	if cacheSizeRaw == "" {
		return Config{}, errors.New("no value found for CACHE_SIZE environment variable")
	}

	cacheSize, err := strconv.Atoi(cacheSizeRaw)
	if err != nil {
		return Config{}, errors.New("CACHE_SIZE must be an integer")
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
