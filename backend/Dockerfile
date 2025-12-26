FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/web

FROM alpine:3.20

WORKDIR /app

RUN adduser -D -g "" appuser

COPY --from=builder /app/server /app/server
COPY --from=builder /app/internal/migrations/json /app/internal/migrations/json
COPY --from=builder /app/api/openapi.yaml /app/api/openapi.yaml

USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/server"]
