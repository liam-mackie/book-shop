# Golang build
FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bookshop ./cmd/bookshop

EXPOSE 8080

ENTRYPOINT ["/app/bookshop"]