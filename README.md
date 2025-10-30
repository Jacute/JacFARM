# JacFARM

[![tests](https://github.com/Jacute/JacFARM/actions/workflows/tests.yml/badge.svg)](https://github.com/Jacute/JacFARM/actions/workflows/tests.yml)
![GitHub Release](https://img.shields.io/github/v/release/Jacute/JacFARM)


[![Coverage Status](https://coveralls.io/repos/github/Jacute/JacFARM/badge.svg?branch=tests-jacfarm-api)](https://coveralls.io/github/Jacute/JacFARM)
[![codecov](https://codecov.io/gh/Jacute/JacFARM/branch/tests-jacfarm-api/graph/badge.svg)](https://app.codecov.io/gh/Jacute/JacFARM)

Exploit farm for attack-defense CTF competition

## Components

### Arch Diagram

![](./docs/img/diagram.jpg)

### Server

- **Exploit Runner** - a worker that launches exploits on all teams. [More details](./docs/exploit_runner/exploit_runner.md)
- **Flag Sender** - a worker that sends flags to jury using *Plugins*. [More details](./docs/flag_sender/flag_sender.md)
- **JacFARM API** - API for frontend and cli start_exploit.py.
- **Config Loader** - loads config into db from config.yml on start. Next configuration editing is available through the frontend.

### Client

- **start_exploit.py** - python cli tool for starting exploits on local machine (TODO)
- **Frontend**
![](./docs/img/frontend.png)