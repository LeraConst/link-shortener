package tests

import (
	"testing"
	"github.com/LeraConst/link-shortener/internal/service"
)

func TestGenerateShortURL(t *testing.T) {
	url := "https://example.com"
	shortURL := service.GenerateShortURL(url)

	// Проверяем длину короткой ссылки
	if len(shortURL) != 10 {
		t.Errorf("Expected short URL of length 10, but got %d", len(shortURL))
	}

	// Тестируем, что ссылка будет всегда одинаковой для одного и того же URL
	shortURL2 := service.GenerateShortURL(url)
	if shortURL != shortURL2 {
		t.Errorf("Generated short URLs do not match: %s != %s", shortURL, shortURL2)
	}
}
