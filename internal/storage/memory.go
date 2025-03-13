package storage

import "sync"

// MemoryStorage — это структура, реализующая хранение ссылок в памяти.
// Она использует map для хранения "сокращённый URL - оригинальный URL".
// Доступ к данным защищён мьютексами, чтобы обеспечить безопасность при выполнении несколько обработчиков одновременно
type MemoryStorage struct {
	data map[string]string
	mu   sync.RWMutex
}

// NewMemoryStorage создает новое хранилище в памяти
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]string),
	}
}

// Save сохраняет ссылку в памяти
func (m *MemoryStorage) Save(originalURL, shortURL string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[shortURL] = originalURL
}

// Get возвращает оригинальный URL
func (m *MemoryStorage) Get(shortURL string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	originalURL, exists := m.data[shortURL]
	if !exists {
		return ""
	}
	return originalURL
}

// CheckExists проверяет, есть ли уже такая ссылка
func (m *MemoryStorage) CheckExists(originalURL string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for short, long := range m.data {
		if long == originalURL {
			return short
		}
	}
	return ""
}
