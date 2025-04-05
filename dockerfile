# Этап сборки
FROM golang:alpine AS builder
# sudo docker run -e CONFIG_PATH=/app/config/local.yaml -p 8080:8080  ttk-back
# sudo docker run -it --entrypoint /bin/sh realslimpudge/ttk-front    npx next start
RUN apk add --no-cache bash git
WORKDIR /app
COPY . .
RUN go get github.com/steebchen/prisma-client-go
RUN go run github.com/steebchen/prisma-client-go generate --schema=./storage/prisma/schema.prisma

RUN go mod tidy

RUN go build -ldflags="-s -w" -o /app/main ./cmd/main.go


# Этап запуска
FROM alpine:latest
RUN apk add --no-cache 
WORKDIR /app

# Копируем бинарник и ВСЮ структуру проекта
COPY --from=builder /app/main .
COPY --from=builder /app/.env /app/
COPY --from=builder /app/config ./config
COPY --from=builder /app/prisma ./prisma

EXPOSE 8080

ENTRYPOINT ["/app/main"]
CMD ["--config=/app/config/local.yaml"]  # Используем абсолютный путь
