# Roadmap

## Core

- For now we only support temperature, humidity, battery, button event capabilities. Handle more capabilities.
- auto restart scanner used by registered device, cause actually if we restart core, scanner doesnt restarts

- find a way to unify value format in core, for exemple determine stricts constraint on temperature format (float, int, ...) and avoid scanners to send in an other format.

## Homekit Adapter

- Handle every device possible to create with HAP

- Find a way to retrieve Homekit QR (or code) from front


## Bluetooth Scanner

### Bthome Parser

- Handle every bthome properties


### Switchbot parser

- Handle all switchbot device types (only meter handled for now)


## New Adapters to implement

- Influx db adapter


