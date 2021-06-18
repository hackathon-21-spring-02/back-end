package domain

import "time"

type File struct {
	ID             string    `json:"id"`
	ComposerID     string    `json:"composer_id"`
	FavoriteCount  uint32     `json:"favorite_count"`
	IsFavoriteByMe bool      `json:"is_favorite_by_me"`
	CreatedAt      time.Time `json:"created_at"`
}

type Composers struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	PostCount int    `json:"post_count"`
}