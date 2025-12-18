# go-home

A lightweight bridge for home devices written in Go.


## 🚀 Getting Started

The easiest way to run `go-home` is using Docker Compose. The system consists of a Core service and optional plugins (Bluetooth, HomeKit, etc.).


### 1. Core Service (Required)

```yaml
services:
  mqtt:
    mem_limit: 50m
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
    mem_limit: 50m
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
      - SESSION_SECRET=<your_secret_here>
      
```

Mosquitto config file (`./config/mosquitto.conf`):
```
listener 1883

allow_anonymous true

persistence true
persistence_location /mosquitto/data/
log_dest file /mosquitto/log/mosquitto.log
```


### 2. Add Plugins (Optional)

Add these services to your docker-compose.yml to enable specific features.

#### Bluetooth scanner plugin

Scans for BLE devices. Requires host network access.

```yaml
gohome-bluetooth:
    mem_limit: 50m
    image: ghcr.io/bastien2203/go-home-bluetooth-scanner:latest
    container_name: gohome-bluetooth
    network_mode: host
    privileged: true
    depends_on:
      - mqtt
      - gohome-core
    volumes:
      - /run/dbus:/run/dbus:ro
    restart: unless-stopped
    environment:
      - BROKER_URL=tcp://localhost:1883
      - ENV=production
      - DBUS_SYSTEM_BUS_ADDRESS=unix:path=/run/dbus/system_bus_socket
```

#### HomeKit adapter plugin

Expose devices to Apple Homekit. Requires host network access (for mDNS).

```yaml
  gohome-homekit:
    mem_limit: 50m
    image: ghcr.io/bastien2203/go-home-homekit-adapter:latest
    container_name: gohome-homekit
    network_mode: host
    privileged: true
    depends_on:
      - mqtt
      - gohome-core
    volumes:
      - ./homekit_data:/homekit_data
    restart: unless-stopped
    environment:
      - BROKER_URL=tcp://localhost:1883
      - ENV=production
      - INTERNET_INTERFACE=wlan0 # change to your internet-facing interface
```


## ⚙️ Configuration

| Variable | Description | Example |
|----------|-------------|---------|
| BROKER_URL | MQTT broker URL | `tcp://mqtt:1883` (Bridge) or `tcp://localhost:1883` (Host) |
| SQLITE_DB_PATH | Path to database file | `./data/gohome.db` |
| API_PORT | Port for the Core API | `9880` |
| SESSION_SECRET | Secret for signing cookies | `random_string` |
| ENV | Environment mode | `production` or `dev` |
| INTERNET_INTERFACE | Network interface for mDNS | `wlan0` or `eth0` |

## TODO
- auto restart scanner used by registered device
- homekit qr code
- influx db adapter



