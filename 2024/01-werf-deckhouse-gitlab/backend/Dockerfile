# Используем многоступенчатую сборку образа (multi-stage build)
# Образ, в котором будет собираться проект
FROM golang:1.18-alpine AS build
# Устанавливаем curl и tar.
RUN apk add curl tar
# Копируем исходники приложения
COPY backend/. /app
WORKDIR /app
# Скачиваем утилиту migrate и распаковываем полученный архив.
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
# Запускаем загрузку нужных пакетов.
RUN go mod download
# Запускаем сборку приложения.
RUN go build -o /goapp cmd/main.go

# Образ, который будет разворачиваться в кластере.
FROM alpine:latest
WORKDIR /
# Копируем из сборочного образа исполняемый файл проекта.
COPY --from=build /goapp /goapp
# Копируем из сборочного образа распакованный файл утилиты migrate и схемы миграции.
COPY --from=build /app/migrate /migrations/migrate
COPY backend/db/migrations /migrations/schemes
EXPOSE 8080
ENTRYPOINT ["/goapp"]
