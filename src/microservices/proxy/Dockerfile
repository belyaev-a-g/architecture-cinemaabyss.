FROM golang:1.23-alpine AS builder

WORKDIR /app

# Копируем файлы go mod и sum
COPY go.mod ./
COPY go.sum ./

# PROXY
#ENV HTTP_PROXY=http://172.16.0.250:4343
#ENV HTTPS_PROXY=http://172.16.0.250:4343

# Загружаем все зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o proxy-service .

# Используем меньший образ для финального приложения
FROM alpine:latest

# PROXY
#ENV HTTP_PROXY=http://172.16.0.250:4343
#ENV HTTPS_PROXY=http://172.16.0.250:4343

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарный файл из builder
COPY --from=builder /app/proxy-service .

# Открываем порт приложения
EXPOSE 8000

# PROXY
#ENV HTTP_PROXY=
#ENV HTTPS_PROXY=

# Команда для запуска исполняемого файла
CMD ["./proxy-service"]
