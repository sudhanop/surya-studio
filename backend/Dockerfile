FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata wget
WORKDIR /app
COPY --from=builder /server .
RUN addgroup -S app && adduser -S -G app app && mkdir -p /app/uploads /app/data && chown -R app:app /app
ENV PORT=8080
EXPOSE 8080
USER app
CMD ["./server"]
