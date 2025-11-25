FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk update && apk add --no-cache \
    gcc \
    bluez \
    bluez-deprecated \
    bluez-libs \
    bluez-dev \
    libc-dev \
    linux-headers \
    dbus


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /app/core cmd/.

RUN mkdir -p /app/bin

RUN for dir in /app/cmd/native-plugins/*; do \
    name=$(basename "$dir"); \
    if [ -d "$dir" ]; then \
        echo "Building $name..."; \
        CGO_ENABLED=1 GOOS=linux go build -o /app/bin/"$name" "$dir"/. ; \
    fi \
done


FROM alpine:latest

RUN apk update && apk add --no-cache \
    ca-certificates \
    bluez \
    bluez-deprecated \
    bluez-libs \
    bluez-dev \
    libc-dev \
    linux-headers \
    dbus \
    dbus-glib

RUN mkdir -p /run/dbus

WORKDIR /app


COPY --from=builder /app/core ./core
COPY --from=builder /app/bin ./bin

RUN chmod +x ./core


EXPOSE 8080

ENTRYPOINT ["./core"]