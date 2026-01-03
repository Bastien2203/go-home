# :material-apple: HomeKit Adapter

The HomeKit Adapter acts as a bridge, exposing your GoHome devices to Apple devices (iPhone, iPad, Apple TV).

## Configuration

Add this service to your `docker-compose.yml`.

!!! note "mDNS Requirement" HomeKit relies on mDNS (Bonjour) for discovery. This requires network_mode: host.

```yaml
gohome-homekit:
    image: ghcr.io/bastien2203/go-home-homekit-adapter:latest
    container_name: gohome-homekit
    mem_limit: 50m
    network_mode: host
    privileged: true
    depends_on:
      - mqtt
      - gohome-core
    volumes:
      - ./homekit_data:/homekit_data # (1)!
    restart: unless-stopped
    environment:
      - BROKER_URL=tcp://localhost:1883
      - ENV=production
      - INTERNET_INTERFACE=wlan0 # (2)!
```

1. Persists pairing data. If you lose this folder, you will have to re-pair everything.

2. Crucial: Set this to your actual network interface name (e.g., `eth0`, `wlan0`, `enp3s0`).

## Capabilities

Currently mapped capabilities between GoHome and HomeKit:

- :material-gesture-tap-button: CapabilityButtonEvent
- :material-battery: CapabilityBattery
- :material-water-percent: CapabilityHumidity
- :material-thermometer: CapabilityTemperature

!!! tip "Pairing" Look at the container logs to find the pairing code or QR code when starting the service for the first time. bash docker logs gohome-homekit