package domain

import "context"

type URLRepository interface {
	Create(ctx context.Context, url *URL) error
	GetBySlug(ctx context.Context, slug string) (*URL, error)
	Delete(ctx context.Context, slug string) error
	IncrementClickCount(ctx context.Context, slug string) error
}

type ClickRepository interface {
	Record(ctx context.Context, click *Click) error
}

type URLService interface {
	Shorten(ctx context.Context, originalURL, customSlug string) (*URL, error)
	Get(ctx context.Context, slug string) (*URL, error)
	Resolve(ctx context.Context, slug string) (*URL, error)
	Delete(ctx context.Context, slug string) error
	SuggestSlugs(ctx context.Context, originalURL string) (<-chan string, error)
}
