package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/LeraConst/link-shortener/internal/storage"
	"github.com/LeraConst/link-shortener/internal/service"
)

func TestResolveHandler(t *testing.T) {
	store := storage.NewMemoryStorage()
	originalURL := "https://example.com"
	shortURL := "AbcD123_Xy"

	// Сохраняем ссылку
	store.Save(originalURL, shortURL)

	// Эмулируем запрос на редирект
	req := httptest.NewRequest(http.MethodGet, "/"+shortURL, nil)
	rec := httptest.NewRecorder()

	// Запускаем обработчик
	handler := service.ResolveHandler(store)
	handler.ServeHTTP(rec, req)

	// Проверяем, что был редирект
	if rec.Code != http.StatusFound {
		t.Errorf("Expected status code 302, but got %d", rec.Code)
	}

	// Проверяем, что в Location заголовке правильный URL
	location := rec.Header().Get("Location")
	if location != originalURL {
		t.Errorf("Expected Location header %s, but got %s", originalURL, location)
	}
}
