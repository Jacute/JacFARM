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

#### Plugins

**Plugin** - is a function in a farm that sends flags to the jury system.

⚠️ Farm contains only one plugin for jury [ForcAD](https://github.com/pomo-mondreganto/ForcAD). If you write plugins for other jury systems, you can create a pull request to add them into repository.

[Plugin example for ForcAD](./workers/flag_sender/plugins/forcad_http/client.go)

### Client

- **start_exploit.py** - python cli tool for starting exploits on local machine (TODO)
- **Frontend** - ui for
  - view flags
  - add exploits of different types (python, binary) via '+' button
  - delete or update exploits by right button
  - add teams
  - update farm config (config can be updated in real-time, except for the flag sending plugin option)
  - view logs of flag_sender and exploit_runner services

![](./docs/img/frontend.png)

#### Examples of exploits
- [binary](docs/exploit_runner/exploit_examples/binary/)
- [python (one file)](docs/exploit_runner/exploit_examples/python_one_file/)
- [python (zip)](docs/exploit_runner/exploit_examples/python_zip/)
