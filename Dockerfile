FROM node:22-bookworm-slim AS node-builder

WORKDIR /app
COPY front/package.json front/package-lock.json ./
RUN npm install

ENV VITE_APP_ENV=production
COPY front/ .
RUN npm run build

FROM --platform=$BUILDPLATFORM golang:1.25-bookworm AS go-builder

COPY --from=tonistiigi/xx:master / /

WORKDIR /app

ARG TARGETPLATFORM

RUN apt-get update && apt-get install -y clang lld
RUN xx-apt-get install -y \
    gcc \
    libsqlite3-dev \
    libc6-dev \
    pkg-config 


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 xx-go build -o /app/core .

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app

COPY --from=node-builder /app/dist ./dist
COPY --from=go-builder /app/core ./core

RUN chmod +x ./core

EXPOSE 8080

ENTRYPOINT ["./core"]