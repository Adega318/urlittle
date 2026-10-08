package store

import (
	"context"
	"crypto/rand"
	"errors"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type Store struct {
	cache *lru.Cache[string, string]
	db    *pgxpool.Pool
}

const (
	URLIDSize     = 6
	maxIDAttempts = 3
	uniqueCode    = "23505"
)

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
	for attempt := 0; attempt < maxIDAttempts; attempt++ {
		id := rand.Text()[:URLIDSize]

		_, err := s.db.Exec(
			ctx,
			"INSERT INTO urls (id, url) VALUES ($1, $2)",
			id,
			value,
		)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == uniqueCode {
				continue
			}
			return "", err
		}

		s.cache.Add(id, value)

		return id, nil
	}

	return "", errors.New("could not generate a unique id")
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
				return "", ErrNotFound
			}

			return "", err
		}

		s.cache.Add(key, value)
	}

	return value, nil
}
