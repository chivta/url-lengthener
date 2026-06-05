FROM golang:1.26.2-alpine AS dev
WORKDIR /app
RUN go install github.com/air-verse/air@latest
CMD ["air"]

FROM golang:1.26.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /api ./cmd/server/main.go

FROM gcr.io/distroless/static:nonroot AS prod
COPY --from=builder /api /api
USER nonroot:nonroot
ENTRYPOINT ["/api"]
