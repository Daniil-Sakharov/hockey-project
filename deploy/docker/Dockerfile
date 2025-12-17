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
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/parser ./cmd/parser
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/stats_parser ./cmd/stats_parser
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/fhspb_parser ./cmd/fhspb_parser
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/fhspb_stats_parser ./cmd/fhspb_stats_parser

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /bin/bot /app/bot
COPY --from=builder /bin/migrate /app/migrate
COPY --from=builder /bin/parser /app/parser
COPY --from=builder /bin/stats_parser /app/stats_parser
COPY --from=builder /bin/fhspb_parser /app/fhspb_parser
COPY --from=builder /bin/fhspb_stats_parser /app/fhspb_stats_parser
COPY --from=builder /app/migrations /app/migrations

ENV TZ=Europe/Moscow
WORKDIR /app

CMD ["/app/bot"]
