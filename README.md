# JacFARM

<p align="center">
  <picture>
    <img alt="JacFARM Logo" src="https://raw.githubusercontent.com/Jacute/JacFARM/refs/heads/master/docs/logo.png" width="300">
  </picture>
</p>
<p align="center">
  <strong>Exploit farm for attack-defense CTF competition</strong>
</p>
<p align="center">
  <a href="#quick-start">Quick Start</a> •
  <a href="#features">Features</a> •
  <a href="#components">Components</a>
</p>
<p align="center">
  <a href="https://github.com/Jacute/JacFARM/actions"><img src="https://github.com/Jacute/JacFARM/actions/workflows/tests.yml/badge.svg" alt="CI Status"></a>
  <a href="https://codecov.io/gh/ollelogdahl/concord"><img alt="Codecov" src="https://codecov.io/gh/Jacute/JacFARM/master/master/graph/badge.svg"></a>
  <a href="#"><img alt="Coveralls" src="https://coveralls.io/repos/github/Jacute/JacFARM/badge.svg?branch=master"></a>
  <a href="https://github.com/Jacute/JacFARM/releases"><img alt="Release" src="https://img.shields.io/github/v/release/Jacute/JacFARM"></a>
</p>

## Quick start

### Dependencies

- Docker
- Docker Compose
- Make

### Start

1. Configure *config.yml* for your competition. A detailed description of the quick configuration is [here](./docs/config.md)

2. Start the farm
```bash
make up
```

Credentials for basic auth and the token for sending flags via start_exploit.py will be printed to stdout.

3. After the game ends, turn off the farm and clean the database and queue
```bash
make down
make clean-all
```

## Features

- Uploading exploits in ui
- Real-time configuration farm options like number of concurrently running exploits, the size of the flag sending batch, team ip addresses, etc
- The ability to [change the plugin for sending flags to jury](./docs/flag_sender/flag_sender.md).
- There are already two sending plugins: [forcad_http](./workers/flag_sender/plugins/forcad_http/client.go) and [saarctf_tcp](./workers/flag_sender/plugins/saarctf_tcp/client.go).
- Different [exploit types](./docs/exploit_runner/exploit_runner.md):
  - Python (one file)
  - Python (zip)
  - Bash script
  - Binary
- View logs of running exploits and sending flags on ui
- Configuring vulnboxes ip addresses using [various methods](./docs/config.md)

## Components

### Client

- **Frontend** - ui for
  - viewing flags with any filters
  - adding exploits of different types via '+' button
  - deleting or updating exploits by right mouse button
  - adding teams
  - updating farm config
  - viewing logs

![](./docs/img/frontend.png)

- **start_exploit.py** - python cli tool for starting exploits on local machine (TODO)

### Server

- **Exploit Runner** - a worker that launches exploits on all teams. [More details](./docs/exploit_runner/exploit_runner.md)
- **Flag Sender** - a worker that sends flags to jury using *Plugins*. [More details](./docs/flag_sender/flag_sender.md)
- **JacFARM API** - API for frontend and cli start_exploit.py.
- **Config Loader** - loads config into db from config.yml on start. Next configuration editing is available through the frontend.

#### Plugins

**Plugin** - is a function in a farm that sends flags to the jury system.

⚠️ Farm contains two plugins for [ForcAD](https://github.com/pomo-mondreganto/ForcAD) and saarCTF jury systems. If you write plugins for other jury systems, you can create a pull request to add them into repository.

[Plugin example for ForcAD](./workers/flag_sender/plugins/forcad_http/client.go)

### Arch Diagram

![](./docs/img/diagram.jpg)