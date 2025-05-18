# Этап сборки (builder)
FROM golang:1.24-alpine AS builder

# Установка зависимостей для сборки
RUN apk add --no-cache git make

# Рабочая директория
WORKDIR /app

# Копируем файлы модулей для кэширования
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
# Флаги:
# - CGO_ENABLED=0 - отключаем CGO для полностью статического бинарника
# - ldflags="-s -w" - уменьшаем размер бинарника
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/main .

# Этап запуска (runtime)
FROM alpine:3.18

# Устанавливаем tzdata для работы с временными зонами
RUN apk add --no-cache tzdata ca-certificates

# Рабочая директория
WORKDIR /app

# Копируем бинарник из этапа сборки
COPY --from=builder /app/main /app/main
# Копируем конфигурационные файлы
COPY configs /app/configs

# Переменные окружения по умолчанию
ENV APP_ENV=production \
    TZ=UTC

# Открываем порт (если приложение использует)
EXPOSE 8080

# Команда запуска
CMD ["/app/main"]