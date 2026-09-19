package service

import (
	"url_shortener/internal/generator"
	"url_shortener/internal/storage"
)
type URLService struct {
	store *storage.MemoryStore // pointer to our storage, so everyone shares the same data
}

func NewURLService(store *storage.MemoryStore) *URLService {
	return &URLService{
		store: store,
	}
}

func (s *URLService) CreateShortURL(longURL string) string {
	code := generator.Generate(6) 
	s.store.Save(code, longURL)   
	return code
}

func (s *URLService) GetOriginalURL(code string) (string, bool) {
	return s.store.Get(code)
}