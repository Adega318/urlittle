package store

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/Adega318/urlittle/internal"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	cache *lru.Cache[string, string]
	db    *pgxpool.Pool
}

func NewStore(ctx context.Context, connString string, cacheSize int) (*Store, error) {
	cache, err := lru.New[string, string](cacheSize)
	if err != nil {
		return nil, err
	}

	db, err := newDb(ctx, connString)
	if err != nil {
		return nil, err
	}

	return &Store{
		cache: cache,
		db:    db,
	}, nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) Add(ctx context.Context, value string) (string, error) {
	id, err := generateID()
	if err != nil {
		return "", err
	}

	value, err = internal.NormalizeURL(value)
	if err != nil {
		return "", err
	}

	_, err = s.db.Exec(
		ctx,
		"INSERT INTO urls (id, url) VALUES ($1, $2)",
		id,
		value,
	)
	if err != nil {
		return "", err
	}

	s.cache.Add(id, value)

	return id, nil
}

func (s *Store) Get(ctx context.Context, key string) (string, error) {
	value, ok := s.cache.Get(key)
	if !ok {
		err := s.db.QueryRow(
			ctx,
			"SELECT url FROM urls WHERE id = $1",
			key,
		).Scan(&value)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return "", errors.New("not found")
			}

			return "", err
		}

		s.cache.Add(key, value)
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
