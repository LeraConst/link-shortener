package storage

import (
	"database/sql"
	"log"
)

type Storage interface {
	Save(originalURL, shortURL string)
	Get(shortURL string) string
	CheckExists(originalURL string) string
}

type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage создает новое подключение к PostgreSQL
func NewPostgresStorage(connStr string) *PostgresStorage {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка при открытии БД:", err)
	}

	// Проверяем соединение
	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка при подключении к БД:", err)
	}
	// Создаём таблицу, если её нет
	CreateTable(db)
	return &PostgresStorage{db: db}
}

// Save сохраняет оригинальный URL и возвращает короткую ссылку
func (p *PostgresStorage) Save(originalURL, shortURL string) {
	_, err := p.db.Exec("INSERT INTO links (short_url, original_url) VALUES ($1, $2)", shortURL, originalURL)
	if err != nil {
		log.Fatal("Ошибка метода Save:", err)
	}
}

// Get возвращает оригинальный URL по короткой ссылке
func (p *PostgresStorage) Get(shortURL string) string {
	var originalURL string
	err := p.db.QueryRow("SELECT original_url FROM links WHERE short_url = $1", shortURL).Scan(&originalURL)
	if err != nil {
		log.Fatal("Ошибка метода Get:", err)
	}
	return originalURL
}

// CheckExists проверяет, существует ли уже такая ссылка
func (p *PostgresStorage) CheckExists(originalURL string) string {
	var shortURL string
	err := p.db.QueryRow("SELECT short_url FROM links WHERE original_url = $1", originalURL).Scan(&shortURL)
	if err == sql.ErrNoRows {
		return ""
	} else if err != nil {
		log.Fatal("Ошибка метода CheckExists:", err)
	}
	return shortURL
}

// CreateTable создаёт таблицу links, если её нет
func CreateTable(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS links (
			id INTEGER PRIMARY KEY,
			short_url VARCHAR(10) UNIQUE NOT NULL,
			original_url TEXT UNIQUE NOT NULL
		);
		`)
	if err != nil {
		log.Fatal("Ошибка при создании таблицы:", err)
	}
	log.Println("БД и таблица links успешно созданы.")
}
