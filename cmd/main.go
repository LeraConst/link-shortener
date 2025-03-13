package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/LeraConst/link-shortener/internal/service"
	"github.com/LeraConst/link-shortener/internal/storage"
	_ "github.com/lib/pq" // Подключение драйвера PostgreSQL
)

var store storage.Storage // интерфейс для работы с разными хранилищами

func main() {
	// Читаем значение из переменных окружения, если они есть
	defaultStorage := os.Getenv("STORAGE_TYPE")
	if defaultStorage == "" {
		defaultStorage = "memory" // Значение по умолчанию
	}
	defaultDbConn := os.Getenv("DB_CONN_STR")
	if defaultDbConn == "" {
		defaultStorage = "postgres://postgres:password@localhost/testdb?sslmode=disable"
	}

	// Флаг для выбора хранилища
	storageType := flag.String("storage", defaultStorage, "Тип хранилища: memory или postgres")
	dbConnStr := flag.String("db", defaultDbConn, "Строка подключения к PostgreSQL")
	flag.Parse()

	// Выбираем хранилище
	switch *storageType {
	case "postgres":
		fmt.Println("Используем хранилище:", *storageType)
		store = storage.NewPostgresStorage(*dbConnStr)
	case "memory":
		fmt.Println("Используем хранилище:", *storageType)
		store = storage.NewMemoryStorage()
	default:
		log.Fatal("Неизвестный тип хранилища")
	}

	// Создаем стандартный маршрутизатор
	mux := http.NewServeMux()

	// Регистрируем обработчики
	mux.HandleFunc("/shorten", service.ShortenHandler(store)) // POST /shorten
	mux.HandleFunc("/", service.ResolveHandler(store))        // GET /{short_url}

	// Запускаем сервер
	fmt.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
