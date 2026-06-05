package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/arvlas/url-shortener/api/internal/domain"
)

type clickRepo struct {
	pool *pgxpool.Pool
}

func NewClickRepo(pool *pgxpool.Pool) domain.ClickRepository {
	return &clickRepo{pool: pool}
}

func (r *clickRepo) Record(ctx context.Context, click *domain.Click) error {
	const q = `
		INSERT INTO clicks (url_id, ip_hash, user_agent, country)
		VALUES ($1, $2, $3, $4)
		RETURNING id, clicked_at`
	return r.pool.QueryRow(ctx, q,
		click.URLID, click.IPHash, click.UserAgent, click.Country,
	).Scan(&click.ID, &click.ClickedAt)
}
