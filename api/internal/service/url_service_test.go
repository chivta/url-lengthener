package service_test

import (
	"context"
	"testing"

	"github.com/arvlas/url-shortener/api/internal/domain"
	"github.com/arvlas/url-shortener/api/internal/service"
)

type mockURLRepo struct {
	createFn              func(ctx context.Context, url *domain.URL) error
	getBySlugFn           func(ctx context.Context, slug string) (*domain.URL, error)
	deleteFn              func(ctx context.Context, slug string) error
	incrementClickCountFn func(ctx context.Context, slug string) error
}

func (m *mockURLRepo) Create(ctx context.Context, url *domain.URL) error {
	if m.createFn != nil {
		return m.createFn(ctx, url)
	}
	return nil
}

func (m *mockURLRepo) GetBySlug(ctx context.Context, slug string) (*domain.URL, error) {
	if m.getBySlugFn != nil {
		return m.getBySlugFn(ctx, slug)
	}
	return nil, nil
}

func (m *mockURLRepo) Delete(ctx context.Context, slug string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, slug)
	}
	return nil
}

func (m *mockURLRepo) IncrementClickCount(ctx context.Context, slug string) error {
	if m.incrementClickCountFn != nil {
		return m.incrementClickCountFn(ctx, slug)
	}
	return nil
}

type mockClickRepo struct {
	recordFn func(ctx context.Context, click *domain.Click) error
}

func (m *mockClickRepo) Record(ctx context.Context, click *domain.Click) error {
	if m.recordFn != nil {
		return m.recordFn(ctx, click)
	}
	return nil
}

func TestShorten_Success(t *testing.T) {
	urlRepo := &mockURLRepo{}
	svc := service.NewURLService(urlRepo, &mockClickRepo{}, nil)

	u, err := svc.Shorten(context.Background(), "https://example.com", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if u.Slug == "" {
		t.Fatal("expected non-empty slug")
	}
}

func TestShorten_InvalidURL(t *testing.T) {
	svc := service.NewURLService(&mockURLRepo{}, &mockClickRepo{}, nil)

	_, err := svc.Shorten(context.Background(), "not-a-url", "")
	if err != domain.ErrInvalidURL {
		t.Fatalf("expected ErrInvalidURL, got %v", err)
	}
}

func TestShorten_InvalidScheme(t *testing.T) {
	svc := service.NewURLService(&mockURLRepo{}, &mockClickRepo{}, nil)

	_, err := svc.Shorten(context.Background(), "ftp://example.com/file.txt", "")
	if err != domain.ErrInvalidURL {
		t.Fatalf("expected ErrInvalidURL, got %v", err)
	}
}

func TestShorten_SlugConflict_RetrySucceeds(t *testing.T) {
	calls := 0
	urlRepo := &mockURLRepo{
		createFn: func(ctx context.Context, url *domain.URL) error {
			calls++
			if calls == 1 {
				return domain.ErrSlugConflict
			}
			return nil
		},
	}
	svc := service.NewURLService(urlRepo, &mockClickRepo{}, nil)

	_, err := svc.Shorten(context.Background(), "https://example.com", "")
	if err != nil {
		t.Fatalf("expected no error after retry, got %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 Create calls, got %d", calls)
	}
}

func TestShorten_CustomSlug_Conflict(t *testing.T) {
	urlRepo := &mockURLRepo{
		createFn: func(ctx context.Context, url *domain.URL) error {
			return domain.ErrSlugConflict
		},
	}
	svc := service.NewURLService(urlRepo, &mockClickRepo{}, nil)

	_, err := svc.Shorten(context.Background(), "https://example.com", "myslug")
	if err != domain.ErrSlugConflict {
		t.Fatalf("expected ErrSlugConflict, got %v", err)
	}
}

func TestGet_NotFound(t *testing.T) {
	urlRepo := &mockURLRepo{
		getBySlugFn: func(ctx context.Context, slug string) (*domain.URL, error) {
			return nil, domain.ErrURLNotFound
		},
	}
	svc := service.NewURLService(urlRepo, &mockClickRepo{}, nil)

	_, err := svc.Get(context.Background(), "missing")
	if err != domain.ErrURLNotFound {
		t.Fatalf("expected ErrURLNotFound, got %v", err)
	}
}

func TestResolve_Success(t *testing.T) {
	want := &domain.URL{Slug: "abc123", OriginalURL: "https://example.com"}
	urlRepo := &mockURLRepo{
		getBySlugFn: func(ctx context.Context, slug string) (*domain.URL, error) {
			return want, nil
		},
	}
	svc := service.NewURLService(urlRepo, &mockClickRepo{}, nil)

	got, err := svc.Resolve(context.Background(), "abc123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.Slug != want.Slug || got.OriginalURL != want.OriginalURL {
		t.Fatalf("expected %+v, got %+v", want, got)
	}
}

func TestResolve_NotFound(t *testing.T) {
	urlRepo := &mockURLRepo{
		getBySlugFn: func(ctx context.Context, slug string) (*domain.URL, error) {
			return nil, domain.ErrURLNotFound
		},
	}
	svc := service.NewURLService(urlRepo, &mockClickRepo{}, nil)

	_, err := svc.Resolve(context.Background(), "gone")
	if err != domain.ErrURLNotFound {
		t.Fatalf("expected ErrURLNotFound, got %v", err)
	}
}

func TestDelete_Success(t *testing.T) {
	svc := service.NewURLService(&mockURLRepo{}, &mockClickRepo{}, nil)

	err := svc.Delete(context.Background(), "abc123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
