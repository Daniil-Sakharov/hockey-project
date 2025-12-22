# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем все бинарники
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/bot ./cmd/bot
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/migrate ./cmd/migrate
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/scheduler ./cmd/scheduler
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/junior_parser ./cmd/junior_parser
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/junior_stats ./cmd/junior_stats
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/fhspb_parser ./cmd/fhspb_parser
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/fhspb_stats ./cmd/fhspb_stats

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /bin/bot /app/bot
COPY --from=builder /bin/migrate /app/migrate
COPY --from=builder /bin/scheduler /app/scheduler
COPY --from=builder /bin/junior_parser /app/junior_parser
COPY --from=builder /bin/junior_stats /app/junior_stats
COPY --from=builder /bin/fhspb_parser /app/fhspb_parser
COPY --from=builder /bin/fhspb_stats /app/fhspb_stats
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/config /app/config

ENV TZ=Europe/Moscow
WORKDIR /app

CMD ["/app/bot"]
