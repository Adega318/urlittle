package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	lru "github.com/hashicorp/golang-lru/v2"
)

type Store struct {
	cache *lru.Cache[string, string]
	db    string // TODO: add db storage
}

func NewStore() *Store {
	cache, _ := lru.New[string, string](100)

	return &Store{cache: cache}
}

func (s *Store) Add(value string) (string, error) {
	id, err := generateID()
	if err != nil {
		return "", err
	}
	value, err = normalizeURL(value)
	if err != nil {
		return "", err
	}

	s.cache.Add(id, value)

	return id, nil
}

func (s *Store) Get(key string) (string, error) {
	value, ok := s.cache.Get(key)
	if !ok {
		return "", errors.New("Not found")
	}

	return value, nil
}

func generateID() (string, error) {
	b := make([]byte, 6) // 48 bits
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
