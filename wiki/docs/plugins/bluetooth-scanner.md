# Bluetooth Scanner

Scans for BLE devices (Xiaomi, SwitchBot, etc.) and pushes data to the Core.

> ⚠️ **Note:** Requires `network_mode: host` and `privileged: true` to access the Bluetooth adapter.

```yaml
  gohome-bluetooth:
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
      # Use localhost because of network_mode: host
      - BROKER_URL=tcp://localhost:1883
      - ENV=production
      - DBUS_SYSTEM_BUS_ADDRESS=unix:path=/run/dbus/system_bus_socket
```


For now just supporting : 
- BHTOME : SensorBattery, SensorHumidity, SensorTemperature, SensorButtonEvent
- SWITCHBOT: Hygrometer and thermometer