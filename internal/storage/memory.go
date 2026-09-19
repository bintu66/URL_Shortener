package storage

type MemoryStore struct {
	data map[string]string }

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]string),
	}
}

func (m *MemoryStore) Save(code string, longURL string) {
	m.data[code] = longURL
}

func (m *MemoryStore) Get(code string) (string, bool) {
	url, exists := m.data[code]
	return url, exists
}