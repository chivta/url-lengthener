package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/rs/zerolog/log"

	"github.com/arvlas/url-lengthener/api/internal/domain"
)

type urlService struct {
	urlRepo   domain.URLRepository
	clickRepo domain.ClickRepository
}

func NewURLService(urlRepo domain.URLRepository, clickRepo domain.ClickRepository) domain.URLService {
	return &urlService{urlRepo: urlRepo, clickRepo: clickRepo}
}

func (s *urlService) Shorten(ctx context.Context, originalURL, customSlug string) (*domain.URL, error) {
	err := validateURL(originalURL)
	if err != nil {
		return nil, domain.ErrInvalidURL
	}

	slug := customSlug
	if slug == "" {
		slug = generateSlug()
	} else if !isValidSlug(slug) {
		return nil, domain.ErrInvalidURL
	}

	u := &domain.URL{Slug: slug, OriginalURL: originalURL}
	err = s.urlRepo.Create(ctx, u)
	if err != nil {
		if errors.Is(err, domain.ErrSlugConflict) && customSlug == "" {
			u.Slug = generateSlug()
			err = s.urlRepo.Create(ctx, u)
			if err != nil {
				return nil, err
			}
			return u, nil
		}
		return nil, err
	}
	return u, nil
}

func (s *urlService) Get(ctx context.Context, slug string) (*domain.URL, error) {
	return s.urlRepo.GetBySlug(ctx, slug)
}

func (s *urlService) Resolve(ctx context.Context, slug string) (*domain.URL, error) {
	u, err := s.urlRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	go func() { //nolint:gosec // intentional: click counting must outlive the request context
		bgCtx := context.Background()
		if err := s.clickRepo.Record(bgCtx, &domain.Click{URLID: u.ID}); err != nil {
			log.Error().Err(err).Str("url_id", u.ID).Msg("click record")
		}
		if err := s.urlRepo.IncrementClickCount(bgCtx, slug); err != nil {
			log.Error().Err(err).Str("slug", slug).Msg("click increment")
		}
	}()
	return u, nil
}

func (s *urlService) Delete(ctx context.Context, slug string) error {
	return s.urlRepo.Delete(ctx, slug)
}

func validateURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("scheme must be http or https")
	}
	if u.Host == "" {
		return fmt.Errorf("missing host")
	}
	return nil
}
