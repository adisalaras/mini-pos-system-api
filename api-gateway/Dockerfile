# builder
FROM golang:1.24-alpine AS builder
WORKDIR /src

# dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build (pastikan main package menghasilkan binary "main")
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /app/main .

# runtime image
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
USER 1000
CMD ["./main"]
