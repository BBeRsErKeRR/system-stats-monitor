# system-stats-monitor
[![Test](https://github.com/BBeRsErKeRR/system-stats-monitor/actions/workflows/pipeline.yml/badge.svg)](https://github.com/BBeRsErKeRR/system-stats-monitor/actions/workflows/pipeline.yml) [![Coverage Status](https://coveralls.io/repos/github/BBeRsErKeRR/system-stats-monitor?branch=develop)](https://coveralls.io/github/BBeRsErKeRR/system-stats-monitor?branch=develop)

## Information

Client/Server GRPC application for monitoring system information.

For more information pleas see [SPECIFICATION](./SPECIFICATION.md)

## Usage

### Daemon

For start daemon you will need:

1. [prepare workspace](#restrictions);
2. Create configuration file with [toml](https://toml.io/en/) format:

```toml
# Logger settings can set stdout/stderr or files
[logger]
level = "DEBUG"
out_paths = ["stdout"]
err_paths = ["stderr"]

[app]
scan_duration = "2s" # How often collect data

[grpc_server]
host = "0.0.0.0"
port = "9080"

# Disable or enable some collectors
[stats]
cpu_enable = true
load_enable = true
network_enable = true
disk_enable = true
network_talkers_enable = true
```

To run daemon use command:

```sh
./ssm --config ./configs/config.toml
```

### Client

For start client you will need create configuration file with [toml](https://toml.io/en/) format:

```toml
# Logger settings can set stdout/stderr or files
[logger]
level = "DEBUG"
out_paths = []
err_paths = []

[app.grpc_client]
termui_enable = true  # Start client with specific UI or print all statistics into stdout
host = "0.0.0.0"
port = "9080"
response_duration = "3s"
wait_duration = "5s"
```

To run client use command:

```sh
./ssm_client --config ./configs/config_client.toml
```

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
- [x] Windows (Limited functionality: cpu only)
