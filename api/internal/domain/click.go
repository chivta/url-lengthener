package domain

import "time"

type Click struct {
	ID        string    `json:"id"`
	URLID     string    `json:"url_id"`
	IPHash    string    `json:"ip_hash"`
	UserAgent string    `json:"user_agent"`
	Country   string    `json:"country"`
	ClickedAt time.Time `json:"clicked_at"`
}
