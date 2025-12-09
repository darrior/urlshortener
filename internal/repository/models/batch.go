package models

type BatchURLs []BatchURLsEntry

type BatchURLsEntry struct {
	ID  string
	URL string
}

type BatchIDs []BatchIDsEntry

type BatchIDsEntry struct {
	ID     string
	UserID string
}
