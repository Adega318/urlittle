package config

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Config struct {
	LogLevel     slog.Level
	Port         string
	Store        StoreConfig
	RateLimiting RateLimitingConfig
}

type StoreConfig struct {
	Ttl       time.Duration
	DbUrl     string
	CacheSize int
}

type RateLimitingConfig struct {
	UserListSize int
	Limit        int
	Burst        int
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
	rateLimitingConfig, err := loadRateLimitingConfig()
	if err != nil {
		return Config{}, err
	}

	return Config{
		LogLevel:     logLevel,
		Port:         port,
		Store:        storeConfig,
		RateLimiting: rateLimitingConfig,
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

func loadRateLimitingConfig() (RateLimitingConfig, error) {
	userListSizeRaw := os.Getenv("RATE_LIMIT_USER_LIST_SIZE")
	if userListSizeRaw == "" {
		return RateLimitingConfig{}, errors.New("no value found for RATE_LIMIT_USER_LIST_SIZE environment variable")
	}
	userListSize, err := strconv.Atoi(userListSizeRaw)
	if err != nil {
		return RateLimitingConfig{}, errors.New("RATE_LIMIT_USER_LIST_SIZE must be an integer")
	}
	if userListSize <= 0 {
		return RateLimitingConfig{}, errors.New("RATE_LIMIT_USER_LIST_SIZE must be positive")
	}

	limitRaw := os.Getenv("RATE_LIMIT_LIMIT")
	if limitRaw == "" {
		return RateLimitingConfig{}, errors.New("no value found for RATE_LIMIT_LIMIT environment variable")
	}
	limit, err := strconv.Atoi(limitRaw)
	if err != nil {
		return RateLimitingConfig{}, errors.New("RATE_LIMIT_LIMIT must be an integer")
	}
	if limit <= 0 {
		return RateLimitingConfig{}, errors.New("RATE_LIMIT_LIMIT must be positive")
	}

	burstRaw := os.Getenv("RATE_LIMIT_BURST")
	if burstRaw == "" {
		return RateLimitingConfig{}, errors.New("no value found for RATE_LIMIT_BURST environment variable")
	}
	burst, err := strconv.Atoi(burstRaw)
	if err != nil {
		return RateLimitingConfig{}, errors.New("RATE_LIMIT_BURST must be an integer")
	}
	if burst <= 0 {
		return RateLimitingConfig{}, errors.New("RATE_LIMIT_BURST must be positive")
	}

	return RateLimitingConfig{UserListSize: userListSize, Limit: limit, Burst: burst}, nil
}
