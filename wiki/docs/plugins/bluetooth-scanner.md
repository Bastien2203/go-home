# Bluetooth Scanner

Scann for bluetooth devices. Requires host network access (for BLE).

## Configuration

Add this service to your docker-compose.yml.

!!! danger "Critical Requirement" This plugin must utilize the host network adapter to access bluetooth hardware.

```yaml
gohome-bluetooth:
    image: ghcr.io/bastien2203/go-home-bluetooth-scanner:latest
    container_name: gohome-bluetooth
    network_mode: host # (1)!
    privileged: true
    depends_on:
      - mqtt
      - gohome-core
    volumes:
      - /run/dbus:/run/dbus:ro # (2)!
    restart: unless-stopped
    environment:
      - BROKER_URL=tcp://localhost:1883 # (3)!
      - ENV=production
      - DBUS_SYSTEM_BUS_ADDRESS=unix:path=/run/dbus/system_bus_socket
```
1. Required to access the physical Bluetooth controller.
2. Maps the host's D-Bus socket to the container, allowing communication with the BlueZ stack.
3. Since we use network_mode: host, we address the broker via localhost, not the docker service name.

## Supported Devices

Currently, the scanner supports the following protocols and capabilities:

### :material-access-point-network: BTHome Standard
Used by many DIY sensors and Xiaomi custom firmwares.

<div class="grid cards" markdown>
- :material-battery: SensorBattery
- :material-water-percent: SensorHumidity
- :material-thermometer: SensorTemperature
- :material-gesture-tap-button: SensorButtonEvent
</div>

### :material-robot: SwitchBot
- Meter (Thermometer/Hygrometer)