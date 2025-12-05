FROM --platform=$BUILDPLATFORM node:22.15-alpine AS node-builder

WORKDIR /app
COPY front/package.json ./
COPY front/package-lock.json ./
RUN npm install

ENV VITE_APP_ENV=production

COPY front/ .
RUN npm run build

FROM golang:1.25-alpine AS go-builder

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

RUN CGO_ENABLED=1 GOOS=linux go build -o /app/core .

RUN mkdir -p /app/bin

RUN for dir in /app/cmd/native-plugins/*; do \
    name=$(basename "$dir"); \
    if [ -d "$dir" ]; then \
        echo "Building $name..."; \
        CGO_ENABLED=1 GOOS=linux go build -o /app/bin/"$name" "$dir"/. ; \
    fi \
done


FROM alpine:3.22.2

RUN apk update

RUN apk add --no-cache libc-dev linux-headers bluez-dev

RUN apk add --no-cache ca-certificates bluez bluez-deprecated bluez-libs dbus dbus-glib

RUN mkdir -p /run/dbus

WORKDIR /app

COPY --from=node-builder /app/dist ./dist
COPY --from=go-builder /app/core ./core
COPY --from=go-builder /app/bin ./bin

RUN chmod +x ./core


EXPOSE 8080

ENTRYPOINT ["./core"]