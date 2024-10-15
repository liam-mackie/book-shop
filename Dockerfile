# Golang build
FROM golang:1.23 AS build

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o bookshop ./cmd/bookshop

FROM scratch
WORKDIR /app/
COPY --from=build /build/bookshop /app/bookshop

EXPOSE 8080

ENTRYPOINT ["/app/bookshop"]