package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/LeraConst/link-shortener/internal/service"
	"github.com/LeraConst/link-shortener/internal/storage"
	_ "github.com/lib/pq" // Подключение драйвера PostgreSQL
)

var store storage.Storage // интерфейс для работы с разными хранилищами

func main() {
	// Флаг для выбора хранилища
	storageType := flag.String("storage", "memory", "Тип хранилища: memory или postgres")
	dbConnStr := flag.String("db", "postgres://user:password@localhost/dbname?sslmode=disable", "Строка подключения к PostgreSQL")
	flag.Parse()

	// Выбираем хранилище
	switch *storageType {
	case "postgres":
		store = storage.NewPostgresStorage(*dbConnStr)
	case "memory":
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
