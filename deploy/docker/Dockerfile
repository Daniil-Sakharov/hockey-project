# === FRONTEND BUILD ===
FROM node:22-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

# === BACKEND BUILD ===
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Собираем только нужные бинарники (без bot, fhspb_*)
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/migrate ./cmd/migrate
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/scheduler ./cmd/scheduler

# === RUNTIME ===
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata

# Backend бинарники
COPY --from=backend-builder /bin/api /app/api
COPY --from=backend-builder /bin/migrate /app/migrate
COPY --from=backend-builder /bin/scheduler /app/scheduler
COPY --from=backend-builder /app/migrations /app/migrations
COPY --from=backend-builder /app/config /app/config

# Frontend статика
COPY --from=frontend-builder /app/frontend/dist /app/frontend/dist

ENV TZ=Europe/Moscow
WORKDIR /app

CMD ["/app/api"]
