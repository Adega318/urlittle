package store

import (
	"context"
	"crypto/rand"
	"errors"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type Store struct {
	cache *lru.Cache[string, Entry]
	db    *pgxpool.Pool
}

type Entry struct {
	URL       string
	ExpiresAt time.Time
}

const (
	URLIDSize     = 6
	maxIDAttempts = 3
	ttl           = time.Minute * 10
	uniqueCode    = "23505"
)

func NewStore(ctx context.Context, connString string, cacheSize int) (*Store, error) {
	cache, err := lru.New[string, Entry](cacheSize)
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
	explires := time.Now().Add(ttl)

	for attempt := 0; attempt < maxIDAttempts; attempt++ {
		id := rand.Text()[:URLIDSize]

		_, err := s.db.Exec(
			ctx,
			"INSERT INTO urls (id, url, expires_at) VALUES ($1, $2, $3)",
			id,
			value,
			explires,
		)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == uniqueCode {
				continue
			}
			return "", err
		}

		s.cache.Add(id, Entry{value, explires})

		return id, nil
	}

	return "", errors.New("could not generate a unique id")
}

func (s *Store) Get(ctx context.Context, key string) (string, error) {
	value, ok := s.cache.Get(key)

	if !ok {
		err := s.db.QueryRow(
			ctx,
			"SELECT url, expires_at FROM urls WHERE id = $1",
			key,
		).Scan(&value.URL, &value.ExpiresAt)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return "", ErrNotFound
			}

			return "", err
		}

		s.cache.Add(key, value)
	}

	return value.URL, nil
}

func (s *Store) ClearExpired(ctx context.Context) error {
	rows, err := s.db.Query(ctx, `
		DELETE FROM urls
		WHERE expires_at <= NOW()
		RETURNING id
	`)
	if err != nil {
		return err
	}

	ids, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return err
	}

	for _, id := range ids {
		s.cache.Remove(id)
	}

	return nil
}
