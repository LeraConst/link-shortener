# 📌 Link Shortener
**Link Shortener** — это сервис для сокращения длинных URL-адресов. Он поддерживает хранение ссылок как в памяти, так и в базе данных PostgreSQL.

---

## 🔥 Функциональность проекта:
- 🔗 Генерация коротких ссылок для длинных URL.
- 🔄 Перенаправление по короткой ссылке на оригинальный URL.
- 📂 Поддержка хранения ссылок в памяти или базе данных PostgreSQL.
- 🏗 REST API для интеграции с клиентами.
- 🐳 Запуск через Docker.
- ✅ Unit-тесты для проверки работоспособности.

---

## 📂 Структура проекта

```plaintext
📦 link-shortener
├── 📂 cmd/              # Точка входа (main.go)
├── 📂internal/
│   ├── 📂service/
│   │   └── handlers.go  # Обработчики API-запросов
│   ├── 📂storage/
│   │   ├── memory.go    # Хранилище в памяти
│   │   ├── postgres.go  # Хранилище в PostgreSQL
├── 📂tests/             # Unit-тесты
├── go.mod               # Зависимости проекта
├── go.sum               # Контрольные суммы зависимостей
├── Dockerfile           # Конфигурация для контейнеризации
├── 📂.github/
│   └── workflows/
│       └── push.yaml    # CI pipeline
```

---

## 🚀 Запуск проекта

### 1️⃣ Локальный запуск (без Docker)
🛠 *Требования: Go 1.22+ и PostgreSQL (если используете базу данных)*

**Установка зависимостей**
```sh
go mod download
```
**Запуск сервера с хранилищем в памяти**
```sh
cd cmd/
go build -o myapp main.go
./myapp --storage=memory
```
**Запуск сервера с PostgreSQL (замените строку подключения)**
```sh
go build -o myapp main.go
./myapp --storage=postgres --db="postgres://user:password@localhost/dbname?sslmode=disable"
```

### 2️⃣ Запуск через Docker
🐳 *Требования: установлен Docker*

**Собираем образ**
```sh
sudo docker build -t link-shortener .
```
**Запускаем контейнер (по умолчанию хранилище в памяти)**
```sh
sudo docker run -p 8080:8080 -e STORAGE_TYPE=memory link-shortener
```
**Запуск с PostgreSQL (замените строку подключения)**
```sh
sudo docker run -p 8080:8080 -e STORAGE_TYPE=postgres -e DB_CONN_STR="postgres://user:password@localhost/dbname?sslmode=disable" link-shortener
```

---

## 🔧 API эндпоинты

### ➕ Создание короткой ссылки (POST /shorten)

**Запрос:**
```sh
POST http://localhost:8080/shorten
```
```json
{
  "url": "https://example.com"
}
```
**Ответ:**
```json
{
    "short_url": "http://localhost:8080/EAaArVRs5q"
}
```

### 🔄 Редирект по короткой ссылке (GET /{short_url})

**Запрос:**
```sh
GET http://localhost:8080/EAaArVRs5q
```
**Ожидаемый результат:** `Редирект на https://example.com`

---

## 📌 Запуск unit-тестов

Запустить тесты можно вручную в терминале:
```sh
go test -v ./tests/
```

---

## 🔄 CI (GitHub Actions)
Проект использует GitHub Actions для автоматического тестирования.

При каждом push в репозиторий автоматически запускаются:

✔️ Проверка кода (go vet)

✔️ Unit-тесты (go test)

Файл: `.github/workflows/push.yaml`

---