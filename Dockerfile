# Golang build
FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.23 AS build

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o bookshop ./cmd/bookshop

FROM --platform=${TARGETPLATFORM:-linux/amd64} scratch
WORKDIR /app/
COPY --from=build /build/bookshop /app/bookshop

EXPOSE 8080

ENTRYPOINT ["/app/bookshop"]