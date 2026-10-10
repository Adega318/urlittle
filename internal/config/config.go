package config

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Config struct {
	LogLevel slog.Level
	Port     string
	Store    StoreConfig
}

type StoreConfig struct {
	Ttl       time.Duration
	DbUrl     string
	CacheSize int
}

func LoadConfig() (Config, error) {
	logLevelRaw := os.Getenv("LOG_LEVEL")
	if logLevelRaw == "" {
		return Config{}, errors.New("no value found for LOG_LEVEL environment variable")
	}
	var logLevel slog.Level
	if err := logLevel.UnmarshalText([]byte(logLevelRaw)); err != nil {
		return Config{}, errors.New("LOG_LEVEL mut be: DEBUG | INFO | WARN | ERROR")
	}

	port := os.Getenv("PORT")
	if port == "" {
		return Config{}, errors.New("no value found for PORT environment variable")
	}

	storeConfig, err := loadStoreConfig()
	if err != nil {
		return Config{}, err
	}

	return Config{
		LogLevel: logLevel,
		Port:     port,
		Store:    storeConfig,
	}, nil
}

func loadStoreConfig() (StoreConfig, error) {
	ttlRaw := os.Getenv("TTL_MINUTES")
	if ttlRaw == "" {
		return StoreConfig{}, errors.New("no value found for TTL_MINUTES environment variable")
	}
	ttl, err := strconv.Atoi(ttlRaw)
	if err != nil {
		return StoreConfig{}, errors.New("TTL_MINUTES must be an integer")
	}

	dburl := os.Getenv("DATABASE_URL")
	if dburl == "" {
		return StoreConfig{}, errors.New("no value found for DATABASE_URL environment variable")
	}

	cacheSizeRaw := os.Getenv("CACHE_SIZE")
	if cacheSizeRaw == "" {
		return StoreConfig{}, errors.New("no value found for CACHE_SIZE environment variable")
	}
	cacheSize, err := strconv.Atoi(cacheSizeRaw)
	if err != nil {
		return StoreConfig{}, errors.New("CACHE_SIZE must be an integer")
	}

	return StoreConfig{Ttl: time.Duration(ttl) * time.Minute, DbUrl: dburl, CacheSize: cacheSize}, nil
}
