# JacFARM

Exploit farm for attack-defense CTF competition

## Components

### Arch Diagram

![](./docs/img/diagram.jpg)

### Server

- **Exploit Runner** - a worker that launches exploits on all teams. [More details](./docs/exploit_runner/exploit_runner.md)
- **Flag Sender** - a worker that sends flags to jury using *Plugins*. [More details](./docs/flag_sender/flag_sender.md)
- **JacFARM API** - API for frontend and cli start_exploit.py.

### Client

- **start_exploit.py** - python cli tool for starting exploits on local machine
- **Frontend**
![](./docs/img/frontend.png)