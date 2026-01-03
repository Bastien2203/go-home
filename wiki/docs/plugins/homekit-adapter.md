# HomeKit Adapter

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


For now only support capabilities : 
- CapabilityButtonEvent
- CapabilityBattery
- CapabilityHumidity
- CapabilityTemperature