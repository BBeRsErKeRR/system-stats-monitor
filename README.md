# system-stats-monitor
[![Test](https://github.com/BBeRsErKeRR/system-stats-monitor/actions/workflows/pipeline.yml/badge.svg)](https://github.com/BBeRsErKeRR/system-stats-monitor/actions/workflows/pipeline.yml) [![Coverage Status](https://coveralls.io/repos/github/BBeRsErKeRR/system-stats-monitor?branch=develop)](https://coveralls.io/github/BBeRsErKeRR/system-stats-monitor?branch=develop)

## Information

Client/Server GRPC application for monitoring system information.

For more information pleas see [SPECIFICATION](./SPECIFICATION.md)

## Usage

## Restrictions

- for linux:
  - installed packages:
    - iostat
    - tcpdump
- sudo NOPASSWD for execute commands:
  - `netstat -lntupe` - for see statistics for all users (optional)
  - `tcpdump -ntq -i any -Q inout -l`
  - `tcpdump -ntq -i any -Q inout -l`

Support os:

- [x] Linux
- [] Window (Limited functionality: cpu only)
