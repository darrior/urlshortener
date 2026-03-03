package api

type UserURLsResponse []UserURLsResponseEntry

type UserURLsResponseEntry struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
