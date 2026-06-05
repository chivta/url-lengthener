package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arvlas/url-shortener/api/internal/config"
	"github.com/arvlas/url-shortener/api/internal/domain"
	"github.com/arvlas/url-shortener/api/internal/handler"
)

type mockURLService struct {
	shortenFn      func(ctx context.Context, originalURL, customSlug string) (*domain.URL, error)
	getFn          func(ctx context.Context, slug string) (*domain.URL, error)
	resolveFn      func(ctx context.Context, slug string) (*domain.URL, error)
	deleteFn       func(ctx context.Context, slug string) error
	suggestSlugsFn func(ctx context.Context, originalURL string) (<-chan string, error)
}

func (m *mockURLService) Shorten(ctx context.Context, originalURL, customSlug string) (*domain.URL, error) {
	if m.shortenFn != nil {
		return m.shortenFn(ctx, originalURL, customSlug)
	}
	return nil, nil
}

func (m *mockURLService) Get(ctx context.Context, slug string) (*domain.URL, error) {
	if m.getFn != nil {
		return m.getFn(ctx, slug)
	}
	return nil, nil
}

func (m *mockURLService) Resolve(ctx context.Context, slug string) (*domain.URL, error) {
	if m.resolveFn != nil {
		return m.resolveFn(ctx, slug)
	}
	return nil, nil
}

func (m *mockURLService) Delete(ctx context.Context, slug string) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, slug)
	}
	return nil
}

func (m *mockURLService) SuggestSlugs(ctx context.Context, originalURL string) (<-chan string, error) {
	if m.suggestSlugsFn != nil {
		return m.suggestSlugsFn(ctx, originalURL)
	}
	ch := make(chan string)
	close(ch)
	return ch, nil
}

func newTestRouter(svc domain.URLService) http.Handler {
	cfg := &config.Config{
		AllowedOrigins: "*",
		Environment:    "test",
		Port:           "8080",
	}
	return handler.NewRouter(cfg, svc)
}

func TestShortenURL_201(t *testing.T) {
	svc := &mockURLService{
		shortenFn: func(ctx context.Context, originalURL, customSlug string) (*domain.URL, error) {
			return &domain.URL{Slug: "abc123", OriginalURL: originalURL}, nil
		},
	}
	body, _ := json.Marshal(map[string]string{"url": "https://example.com"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/urls", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	newTestRouter(svc).ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if resp["slug"] != "abc123" {
		t.Fatalf("expected slug abc123, got %v", resp["slug"])
	}
}

func TestShortenURL_422_MissingBody(t *testing.T) {
	svc := &mockURLService{}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/urls", bytes.NewReader([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	newTestRouter(svc).ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
}

func TestShortenURL_422_InvalidURL(t *testing.T) {
	svc := &mockURLService{
		shortenFn: func(ctx context.Context, originalURL, customSlug string) (*domain.URL, error) {
			return nil, domain.ErrInvalidURL
		},
	}
	body, _ := json.Marshal(map[string]string{"url": "not-valid"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/urls", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	newTestRouter(svc).ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}
	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if resp["code"] != "invalid_url" {
		t.Fatalf("expected code invalid_url, got %v", resp["code"])
	}
}

func TestShortenURL_409_SlugConflict(t *testing.T) {
	svc := &mockURLService{
		shortenFn: func(ctx context.Context, originalURL, customSlug string) (*domain.URL, error) {
			return nil, domain.ErrSlugConflict
		},
	}
	body, _ := json.Marshal(map[string]string{"url": "https://example.com", "custom_slug": "taken"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/urls", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	newTestRouter(svc).ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d", w.Code)
	}
	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if resp["code"] != "slug_conflict" {
		t.Fatalf("expected code slug_conflict, got %v", resp["code"])
	}
}

func TestGetURL_200(t *testing.T) {
	svc := &mockURLService{
		getFn: func(ctx context.Context, slug string) (*domain.URL, error) {
			return &domain.URL{Slug: slug, OriginalURL: "https://example.com"}, nil
		},
	}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/urls/abc", nil)
	w := httptest.NewRecorder()

	newTestRouter(svc).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetURL_404(t *testing.T) {
	svc := &mockURLService{
		getFn: func(ctx context.Context, slug string) (*domain.URL, error) {
			return nil, domain.ErrURLNotFound
		},
	}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/urls/missing", nil)
	w := httptest.NewRecorder()

	newTestRouter(svc).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if resp["code"] != "url_not_found" {
		t.Fatalf("expected code url_not_found, got %v", resp["code"])
	}
}

func TestRedirectURL_301(t *testing.T) {
	svc := &mockURLService{
		resolveFn: func(ctx context.Context, slug string) (*domain.URL, error) {
			return &domain.URL{Slug: slug, OriginalURL: "https://example.com"}, nil
		},
	}
	req := httptest.NewRequest(http.MethodGet, "/abc", nil)
	w := httptest.NewRecorder()

	newTestRouter(svc).ServeHTTP(w, req)

	if w.Code != http.StatusMovedPermanently {
		t.Fatalf("expected 301, got %d", w.Code)
	}
	if loc := w.Header().Get("Location"); loc != "https://example.com" {
		t.Fatalf("expected Location https://example.com, got %q", loc)
	}
}

func TestRedirectURL_404(t *testing.T) {
	svc := &mockURLService{
		resolveFn: func(ctx context.Context, slug string) (*domain.URL, error) {
			return nil, domain.ErrURLNotFound
		},
	}
	req := httptest.NewRequest(http.MethodGet, "/gone", nil)
	w := httptest.NewRecorder()

	newTestRouter(svc).ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestDeleteURL_204(t *testing.T) {
	svc := &mockURLService{}
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/urls/abc", nil)
	w := httptest.NewRecorder()

	newTestRouter(svc).ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestHealthEndpoint(t *testing.T) {
	svc := &mockURLService{}
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	newTestRouter(svc).ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}
