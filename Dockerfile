# ---- Build stage ----
FROM golang:1.26-alpine AS builder

WORKDIR /src

# Leverage Docker layer caching for dependencies.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/holidays-api ./cmd/api

# ---- Runtime stage ----
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /out/holidays-api ./holidays-api
COPY --from=builder /src/docs ./docs

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["./holidays-api"]

