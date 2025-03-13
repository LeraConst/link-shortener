package tests

import (
	"database/sql"
	"log"
	"testing"

	"github.com/LeraConst/link-shortener/internal/storage"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupPostgresTestDB(t *testing.T) *storage.PostgresStorage {
	connStr := "postgres://postgres:password@localhost:5432/dbname?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
		t.FailNow()
	}

	// Удаляем таблицу перед тестами
	_, err = db.Exec("DROP TABLE IF EXISTS links")
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}

	// Создаем таблицу
	_, err = db.Exec(`CREATE TABLE links (
		id SERIAL PRIMARY KEY,
		short_url VARCHAR(10) UNIQUE NOT NULL,
		original_url TEXT UNIQUE NOT NULL
	)`)
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	// Создаем хранилище
	store := storage.NewPostgresStorage(connStr)
	return store
}

func TestPostgresStorage(t *testing.T) {
	store := setupPostgresTestDB(t)

	originalURL := "https://example.com"
	shortURL := "AbcD123_Xy"

	// Тестируем сохранение и получение данных
	store.Save(originalURL, shortURL)

	// Проверяем, что ссылка извлекается корректно
	result := store.Get(shortURL)
	assert.Equal(t, originalURL, result)

	// Проверяем существование ссылки
	existingShort := store.CheckExists(originalURL)
	assert.Equal(t, shortURL, existingShort)
}
