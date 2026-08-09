# icecast2mqtt

Lightweight bridge that polls [Icecast](https://icecast.org) radio servers and publishes their live statistics to MQTT. Includes Home Assistant MQTT auto discovery feature support, so stations show up as MQTT devices with sensors automatically and therefore ready for use in visualizations or automations.

## Overview

Each configured target is polled on its own schedule. On every poll the bridge:

1. Fetches and parses the Icecast status JSON data
2. Publishes server-level metrics to `{base_topic}/{target_topic}`
3. Publishes per-stream metrics to `{base_topic}/{target_topic}/{mountpoint}`
4. Publishes Home Assistant MQTT auto discovery announcements, if `ha_discovery` is enabled for that target

### MQTT metrics

Server-level (`{base_topic}/{target_topic}/`):

| Sub-topic             | Description                            |
|-----------------------|----------------------------------------|
| `version`             | Icecast server version string          |
| `server_start_date`   | Server start time (ISO 8601)           |
| `active_sources`      | Number of currently active mountpoints |

Per-stream (`{base_topic}/{target_topic}/{mountpoint}/`):

| Sub-topic            | Description                    |
|----------------------|--------------------------------|
| `listeners`          | Current listener count         |
| `stream_start_date`  | Stream start time (ISO 8601)   |
| `track`              | Current song, "Artist - Title" |

## Configuration

The service reads `config.json` from the working directory. A static JSON schema
([`config.schema.json`](./config.schema.json)) is linked from the config via the
`$schema` key, so most IDEs (VS Code, JetBrains) provide completion, validation and inline
hints while you edit it. A ready-to-fill example lives in
[`config.json.example`](./config.json.example).

> **MQTT credentials come from environment variables, not from the config file.**

| Environment variable | Required | Description                              |
|----------------------|----------|------------------------------------------|
| `MQTT_BROKER`        | yes      | Broker URL, e.g. `tcp://mqtt.local:1883` |
| `MQTT_USERNAME`      | no       | MQTT username                            |
| `MQTT_PASSWORD`      | no       | MQTT password                            |

## Deployment with Docker

1) Prepare your `config.json` file (copy and edit [`config.json.example`](./config.json.example)):

2) Use this preset for `docker-compose.yml`:

```yaml
services:
  icecast2mqtt:
    image: ghcr.io/SNMetamorph/icecast2mqtt:latest
    container_name: icecast2mqtt
    restart: unless-stopped
    environment:
      - MQTT_BROKER=${MQTT_BROKER:?set MQTT_BROKER in .env}
      - MQTT_USERNAME=${MQTT_USERNAME:-}
      - MQTT_PASSWORD=${MQTT_PASSWORD:-}
    volumes:
      - ./config.json:/app/config.json:ro
```

3) Create a `.env` file next to `docker-compose.yml`:

```sh
MQTT_BROKER=tcp://127.0.0.1:1883
MQTT_USERNAME=your_username
MQTT_PASSWORD=your_password
```

4) Start the service:

```sh
docker compose up -d
```

Your `config.json` is mounted into the container read-only, so editing it on the host and
restarting container is enough to apply changes.

## Home Assistant integration

A short overview of how the bridge interoperates with Home Assistant:

- When `ha_discovery` is set for a target, the bridge publishes MQTT discovery messages under `ha_autodiscovery_prefix` on every poll.
- Each Icecast server becomes a **MQTT device** identified by `device_id` with the friendly name from `device_name`.
- Discovered entities have `sensor` class, the server-level ones are placed in the `diagnostic` entity category. Stream-level metrics are placed as regular sensor entities.

## Building from source

### Prerequisites

- [Go](https://go.dev/dl/) (1.26 or newer)
- [Git](https://git-scm.com/)
- [Task](https://taskfile.dev/) (v3)

### Installing Task utility

```sh
# make sure $(go env GOPATH)/bin is on PATH
go install github.com/go-task/task/v3/cmd/task@latest
```

### Building on Linux and Windows

```sh
task build
```

This compiles `./cmd/app` into `bin/icecast2mqtt` directory.

You could obtain other useful tasks using command `task --list-all`.
