package service

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/LeraConst/link-shortener/internal/storage"
)

// Request структура для входного JSON
type shortenRequest struct {
	URL string `json:"url"`
}

// Response структура для ответа
type shortenResponse struct {
	ShortURL string `json:"short_url"`
}

// ShortenHandler обрабатывает создание короткой ссылки
func ShortenHandler(store storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.Error(res, `{"error": "Метод не поддерживается"}`, http.StatusMethodNotAllowed)
			return
		}

		// Читаем все тело запроса
		body, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(res, `{"error": "Ошибка чтения тела запроса"}`, http.StatusInternalServerError)
			return
		}

		// Декодируем JSON в структуру
		var request shortenRequest
		err = json.Unmarshal(body, &request)
		if err != nil || request.URL == "" {
			http.Error(res, `{"error": "Неверный запрос"}`, http.StatusBadRequest)
			return
		}

		// Проверяем, есть ли уже такая ссылка в БД
		existingShort := store.CheckExists(request.URL)
		if existingShort != "" {
			// Если уже есть короткая ссылка
			// Создаём структуру с короткой ссылкой
			response := shortenResponse{ShortURL: "http://localhost:8080/" + existingShort}

			// Преобразуем структуру в JSON
			responseData, err := json.Marshal(response)
			if err != nil {
				http.Error(res, `{"error": "Ошибка формирования ответа"}`, http.StatusInternalServerError)
				return
			}

			// Записываем JSON в ответ
			res.Header().Set("Content-Type", "application/json")
			if _, err := res.Write(responseData); err != nil {
				http.Error(res, `{"error": "Ошибка записи в ответ"}`, http.StatusInternalServerError)
				return
			}
		}

		// Генерируем новую короткую ссылку
		shortURL := GenerateShortURL(request.URL)

		// Сохраняем в хранилище
		store.Save(request.URL, shortURL)

		// Отправляем ответ с новой короткой ссылкой
		response := shortenResponse{ShortURL: "http://localhost:8080/" + shortURL}
		responseData, err := json.Marshal(response)
		if err != nil {
			http.Error(res, `{"error": "Ошибка формирования ответа"}`, http.StatusInternalServerError)
			return
		}

		// Записываем JSON в ответ
		res.Header().Set("Content-Type", "application/json")
		if _, err := res.Write(responseData); err != nil {
			http.Error(res, `{"error": "Ошибка записи в ответ"}`, http.StatusInternalServerError)
			return
		}
	}
}

// ResolveHandler обрабатывает редирект по короткой ссылке
func ResolveHandler(store storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		// Получаем короткий URL из пути запроса
		shortURL := strings.TrimPrefix(req.URL.Path, "/")
		if shortURL == "" {
			http.Error(res, "Короткая ссылка не найдена", http.StatusBadRequest)
			return
		}

		// Ищем оригинальный URL в БД
		originalURL := store.Get(shortURL)
		if originalURL == "" {
			http.Error(res, "Ссылка не найдена", http.StatusNotFound)
			return
		}

		// Делаем редирект на оригинальный URL
		http.Redirect(res, req, originalURL, http.StatusFound)
	}
}

// GenerateShortURL генерирует короткий код из хеша SHA-256
func GenerateShortURL(url string) string {
	hash := sha256.Sum256([]byte(url))                    // возвращаем хеш в виде массива из 32 байт
	encoded := base64.URLEncoding.EncodeToString(hash[:]) // кодируем в base64
	return strings.TrimRight(encoded, "=")[:10]           // берем 10 символов, убираем '='
}
