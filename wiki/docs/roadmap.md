# :material-map-marker-path: Roadmap

This document outlines the planned features and improvements for the GoHome ecosystem.

## Core

- [ ] Expanded Capabilities: Support more than just temp/humidity (e.g., Lights, Covers, Switches).

- [ ] Auto-Restart Logic: Automatically restart the scanner if a registered device stops responding (currently, scanner restart issues affect core availability).

- [ ] Data Sanitization:

    - [ ] Unify value formats in Core.

    - [ ] Enforce strict constraints (e.g., ensure temperature is always float).

    - [ ] reject malformed payloads from scanners.

## Homekit Adapter

- [ ] Device Parity: Handle every device type possible within the HAP (HomeKit Accessory Protocol) spec.

- [ ] Frontend Integration: Retrieve the HomeKit QR code (or PIN) directly from the Core UI instead of container logs.


## Bluetooth Scanner

### BTHome Parser

- [ ] Handle full BTHome V2 specification (all properties).

### SwitchBot Parser

- [ ] Handle SwitchBot Bot (finger presser).

- [ ] Handle SwitchBot Curtain.

- [ ] Handle SwitchBot Contact Sensor.


## Future Adapters

- [ ] InfluxDB: Adapter to export historical data to InfluxDB for visualization in Grafana.