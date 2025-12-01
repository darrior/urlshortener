package repository

type record struct {
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
}

type urlStorage map[string]record
