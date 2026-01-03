# :material-download: Installation

The recommended way to run `go-home` is using Docker Compose. This ensures that all services (Core, Broker, and Plugins) can communicate easily.

## 1. Prerequisites
* Docker & Docker Compose installed.
* :material-linux: Linux Environment: Highly recommended, especially for plugins requiring network_mode: host (like Bluetooth or HomeKit).


## 2. Configuration & Deployment 

Create a folder for your project and create a docker-compose.yml file.

### Step A: The Docker Compose File
 
Copy the following configuration. It includes the MQTT Broker and the GoHome Core.

```yaml
services:
  # --- MQTT Broker ---
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

  # --- GoHome Core ---
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
      - BROKER_URL=tcp://mqtt:1883 # (1)!
      - SQLITE_DB_PATH=./data/gohome.db
      - API_PORT=9880
      - ENV=production # (2)!
      - SESSION_SECRET=change_me_please # (3)!
      - DEBUG=false
```

1. Inside the Docker network, the hostname is the service name (mqtt). If running outside docker, use localhost.
2. Use dev for non secure, production for stability and security (https, secure cookies, ...).
3. Important: Change this string to something random to secure your session cookies.

### Step B: Mosquitto Configuration

The broker needs a configuration file to allow connections. Create a file at `./config/mosquitto.conf`:


```conf
listener 1883
allow_anonymous true
persistence true
persistence_location /mosquitto/data/
log_dest file /mosquitto/log/mosquitto.log
```

!!! warning "Security Note" allow_anonymous true is useful for local testing. For production usage exposed to the internet, consider setting up a username and password.

## 3. Start the server

Run the stack in detached mode:

```sh
docker-compose up -d
```

Now access dashboard at http://localhost:9880{ .md-button }.