# Installation

The easiest way to run `go-home` is using Docker Compose.

## 1. Prerequisites
* Docker & Docker Compose installed.
* A Linux environment is recommended for Plugins requiring host network access.

## 2. Core Setup (Required)
Create a `docker-compose.yml` file with the MQTT broker and the Core service:

```yaml
services:
  mqtt:
    image: eclipse-mosquitto:2
    container_name: gohome-mqtt
    restart: unless-stopped
    ports:
      - "1883:1883"
      - "9001:9001"
    volumes:
      - ./config:/mosquitto/config
      - /volumes/gohome/mosquitto/data:/mosquitto/data
      - /volumes/gohome/mosquitto/log:/mosquitto/log

  gohome-core:
    image: ghcr.io/bastien2203/go-home-core:latest
    container_name: gohome-core
    ports:
      - "9880:9880"
    depends_on:
      - mqtt
    volumes:
      - ./data:/app/data
    restart: unless-stopped
    environment:
      - BROKER_URL=tcp://mqtt:1883
      - SQLITE_DB_PATH=./data/gohome.db
      - API_PORT=9880
      - ENV=production
      - SESSION_SECRET=change_me_please
```

| Variable | Description | Example |
|----------|-------------|---------|
| BROKER_URL | MQTT broker URL | `tcp://mqtt:1883` (Bridge) or `tcp://localhost:1883` (Host) |
| SQLITE_DB_PATH | Path to database file | `./data/gohome.db` |
| API_PORT | Port for the Core API | `9880` |
| SESSION_SECRET | Secret for signing cookies | `random_string` |
| ENV | Environment mode | `production` or `dev` |
| DEBUG | Debug mode | true or false | 

## 3. Mosquitto Configuration

Create a file at `./config/mosquitto.conf`:

```conf
listener 1883
allow_anonymous true
persistence true
persistence_location /mosquitto/data/
log_dest file /mosquitto/log/mosquitto.log
```


Run docker-compose up -d and access the core at http://localhost:9880.