package repository

type MapRepository struct {
	urls map[string]string
}

var _ Repository = (*MapRepository)(nil)

func NewMapRepository() *MapRepository {
	return &MapRepository{
		urls: map[string]string{},
	}
}

func (r *MapRepository) AddURL(id, url string) error {
	r.urls[id] = url

	return nil
}

func (r *MapRepository) GetURL(id string) (string, error) {
	url, ok := r.urls[id]
	if !ok {
		return "", ErrorNotFound
	}
	return url, nil
}
