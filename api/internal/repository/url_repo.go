package repository

import (
	"context"
	"crypto/sha256"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/arvlas/url-lengthener/api/internal/domain"
)

type urlRepo struct {
	pool *pgxpool.Pool
}

func NewURLRepo(pool *pgxpool.Pool) domain.URLRepository {
	return &urlRepo{pool: pool}
}

// slugHash derives the fixed-size lookup key stored alongside the slug.
// Slugs can be up to 8KB, far past what a btree index can hold, so every
// lookup and uniqueness check goes through this bounded digest instead.
func slugHash(slug string) []byte {
	sum := sha256.Sum256([]byte(slug))
	return sum[:]
}

func (r *urlRepo) Create(ctx context.Context, url *domain.URL) error {
	const q = `
		INSERT INTO urls (slug, slug_hash, original_url, user_id, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, click_count, created_at`

	err := r.pool.QueryRow(ctx, q,
		url.Slug, slugHash(url.Slug), url.OriginalURL, url.UserID, url.ExpiresAt,
	).Scan(&url.ID, &url.ClickCount, &url.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrSlugConflict
		}
		return err
	}
	return nil
}

func (r *urlRepo) GetBySlug(ctx context.Context, slug string) (*domain.URL, error) {
	const q = `
		SELECT id, slug, original_url, user_id, click_count, expires_at, created_at
		FROM urls WHERE slug_hash = $1`

	u := &domain.URL{}
	err := r.pool.QueryRow(ctx, q, slugHash(slug)).Scan(
		&u.ID, &u.Slug, &u.OriginalURL, &u.UserID,
		&u.ClickCount, &u.ExpiresAt, &u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrURLNotFound
		}
		return nil, err
	}
	if u.ExpiresAt != nil && u.ExpiresAt.Before(time.Now()) {
		return nil, domain.ErrURLExpired
	}
	return u, nil
}

func (r *urlRepo) Delete(ctx context.Context, slug string) error {
	const q = `DELETE FROM urls WHERE slug_hash = $1`
	tag, err := r.pool.Exec(ctx, q, slugHash(slug))
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrURLNotFound
	}
	return nil
}

func (r *urlRepo) IncrementClickCount(ctx context.Context, slug string) error {
	const q = `UPDATE urls SET click_count = click_count + 1 WHERE slug_hash = $1`
	_, err := r.pool.Exec(ctx, q, slugHash(slug))
	return err
}
