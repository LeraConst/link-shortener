# Этап сборки
FROM golang:1.22 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для загрузки зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости Go
RUN go mod download

# Копируем все
COPY . .

# Кросс-компиляция Go-программы для Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cmd/myapp cmd/main.go

# Порты
EXPOSE 8080

# Переменные окружения для хранилища
ENV STORAGE_TYPE memory
ENV DB_CONN_STR postgres://postgres:password@localhost/dbname?sslmode=disable

# Команда запуска
CMD ["./cmd/myapp"]
