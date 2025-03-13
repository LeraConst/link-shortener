package tests

import (
	"testing"

	"github.com/LeraConst/link-shortener/internal/storage"
)

func TestMemoryStorage(t *testing.T) {
	store := storage.NewMemoryStorage()

	// Тестируем сохранение и получение данных
	originalURL := "https://example.com"
	shortURL := "AbcD123_Xy"

	// Сохраняем ссылку
	store.Save(originalURL, shortURL)

	// Получаем ссылку
	result := store.Get(shortURL)
	if result != originalURL {
		t.Errorf("Ожидалось %s, получено %s", originalURL, result)
	}

	// Проверяем, что ссылки одинаковые
	existingShort := store.CheckExists(originalURL)
	if existingShort != shortURL {
		t.Errorf("Ожидался короткий URL %s, получен %s", shortURL, existingShort)
	}
}
