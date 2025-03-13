package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/LeraConst/link-shortener/internal/storage"
	"github.com/LeraConst/link-shortener/internal/service"
)

func TestShortenHandler(t *testing.T) {
	store := storage.NewMemoryStorage()

	// Эмулируем запрос с URL
	requestData := map[string]string{"url": "https://example.com"}
	requestBody, _ := json.Marshal(requestData)
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(requestBody))
	rec := httptest.NewRecorder()

	// Запускаем обработчик
	handler := service.ShortenHandler(store)
	handler.ServeHTTP(rec, req)

	// Проверяем, что код ответа 200
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code 200, but got %d", rec.Code)
	}

	// Проверяем, что в ответе есть короткая ссылка
	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response["short_url"] == "" {
		t.Errorf("Expected short_url in response, but got empty")
	}
}
