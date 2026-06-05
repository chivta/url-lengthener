package domain

import "time"

type URL struct {
	ID          string     `json:"id"`
	Slug        string     `json:"slug"`
	OriginalURL string     `json:"original_url"`
	UserID      *string    `json:"user_id,omitempty"`
	ClickCount  int64      `json:"click_count"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}
