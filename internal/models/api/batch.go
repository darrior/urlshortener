package api

type ShortenerBatchRequest []BatchRequestEntry

type BatchRequestEntry struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ShortenerResponse struct {
	Result string `json:"result"`
}

type ShortenerBatchResponse []BatchResponseEntry

type BatchResponseEntry struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
