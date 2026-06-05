package domain

import "errors"

var (
	ErrURLNotFound  = errors.New("url not found")
	ErrSlugConflict = errors.New("slug already exists")
	ErrURLExpired   = errors.New("url has expired")
	ErrInvalidURL   = errors.New("invalid url")
)
