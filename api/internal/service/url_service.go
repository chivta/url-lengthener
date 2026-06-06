package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/anthropics/anthropic-sdk-go"

	"github.com/arvlas/url-lengthener/api/internal/domain"
)

type urlService struct {
	urlRepo   domain.URLRepository
	clickRepo domain.ClickRepository
	claude    *anthropic.Client
}

func NewURLService(urlRepo domain.URLRepository, clickRepo domain.ClickRepository, claude *anthropic.Client) domain.URLService {
	return &urlService{urlRepo: urlRepo, clickRepo: clickRepo, claude: claude}
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
		_ = s.clickRepo.Record(bgCtx, &domain.Click{URLID: u.ID})
		_ = s.urlRepo.IncrementClickCount(bgCtx, slug)
	}()
	return u, nil
}

func (s *urlService) Delete(ctx context.Context, slug string) error {
	return s.urlRepo.Delete(ctx, slug)
}

func (s *urlService) SuggestSlugs(ctx context.Context, originalURL string) (<-chan string, error) {
	ch := make(chan string, 5)
	go func() {
		defer close(ch)
		s.streamSlugs(ctx, originalURL, ch)
	}()
	return ch, nil
}

func (s *urlService) streamSlugs(ctx context.Context, originalURL string, ch chan<- string) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	stream := s.claude.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeHaiku4_5,
		MaxTokens: 256,
		System: []anthropic.TextBlockParam{
			{Text: "You are a URL slug suggester. Call suggest_slugs with exactly 5 memorable, URL-safe slugs (4-8 alphanumeric chars) based on the URL content."},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(fmt.Sprintf("Suggest slugs for: %s", originalURL))),
		},
		Tools: []anthropic.ToolUnionParam{
			{OfTool: &anthropic.ToolParam{
				Name:        "suggest_slugs",
				Description: anthropic.String("Suggest short URL slugs"),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]any{
						"candidates": map[string]any{
							"type": "array",
							"items": map[string]any{
								"type":    "string",
								"pattern": "^[a-zA-Z0-9]{4,8}$",
							},
							"minItems": 5,
							"maxItems": 5,
						},
					},
				},
			}},
		},
	})

	var jsonBuf string
	for stream.Next() {
		event := stream.Current()
		if delta, ok := event.AsAny().(anthropic.ContentBlockDeltaEvent); ok {
			if inputDelta, ok := delta.Delta.AsAny().(anthropic.InputJSONDelta); ok {
				jsonBuf += inputDelta.PartialJSON
			}
		}
	}
	err := stream.Err()
	if err != nil {
		slog.Error("suggest stream", "error", err)
		return
	}

	var result struct {
		Candidates []string `json:"candidates"`
	}
	err = json.Unmarshal([]byte(jsonBuf), &result)
	if err != nil {
		slog.Error("suggest parse", "error", err, "json", jsonBuf)
		return
	}
	for _, slug := range result.Candidates {
		if isValidSlug(slug) {
			select {
			case ch <- slug:
			case <-ctx.Done():
				return
			}
		}
	}
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
