package repository

type record struct {
	OriginalURL string `json:"original_url"`
	UserID      string `json:"user_id"`
	Deleted     bool   `json:"deleted"`
}

type urlStorage map[string]record
