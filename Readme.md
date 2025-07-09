# RabbitMQ Golang Stream Consumer Benchmark

This is a simple Go application designed to benchmark consumer performance for RabbitMQ streams. It connects to one or more streams, consumes messages, and logs consumption statistics (total messages, messages per second, and average entries per message) to the standard output.

## Prerequisites

-   Go (version 1.18 or later)
-   A running RabbitMQ instance with the `stream` plugin enabled.

## Getting Started

### 1. Build the Application

First, build the executable:

```sh
go build .
```

This will create an executable file named `rmq-bench` in the project directory.

### 2. Run the Application

You can run the application with default settings or customize its behavior using command-line flags or environment variables.

To run with default settings (connects to `localhost:5552`, stream `stream1`):

```sh
./rmq-bench
```

Press `Ctrl+C` to gracefully shut down the application.

## Configuration

The application can be configured via command-line flags or environment variables. If both are present, command-line flags take precedence.

| Flag                | Environment Variable    | Description                                       | Default        |
| ------------------- | ----------------------- | ------------------------------------------------- | -------------- |
| `--streams`         | `RMQ_STREAMS`           | Comma-separated list of streams to consume from   | `stream1`      |
| `--host`            | `RMQ_HOST`              | RabbitMQ host                                     | `localhost`    |
| `--vhost`           | `RMQ_VHOST`             | RabbitMQ vhost                                    | `/`            |
| `--port`            | `RMQ_PORT`              | RabbitMQ stream port                              | `5552`         |
| `--user`            | `RMQ_USER`              | RabbitMQ user                                     | `guest`        |
| `--password`        | `RMQ_PASSWORD`          | RabbitMQ password                                 | `guest`        |
| `--max-consumers`   | `RMQ_MAX_CONSUMERS`     | Max consumers per stream                          | `5`            |
| `--timeout`         | `RMQ_TIMEOUT`           | Producer timeout                                  | `10s`          |
| `--initial-credits` | `RMQ_INITIAL_CREDITS`   | Initial credits                                   | `100`          |
| `--log-interval`    | `LOG_INTERVAL`          | Log interval in seconds                           | `5`            |
| `--crc-check`       | `RMQ_CRC_CHECK`         | Enable CRC check for consumers                    | `true`         |

## Examples

### Running with default settings

This command connects to a stream named `stream1` on `localhost:5552`.

```sh
./rmq-bench
```

### Connecting to multiple streams

To consume from `stream-a` and `stream-b`:

```sh
./rmq-bench --streams stream-a,stream-b
```

### Disabling CRC Check

By default, the CRC check is enabled. To disable it:

```sh
./rmq-bench --crc-check=false
```

### Connecting to a different host and port

```sh
./rmq-bench --host my-rabbit.example.com --port 5553
```

### Changing the log interval

To log statistics every 1 second instead of the default 5 seconds:

```sh
./rmq-bench --log-interval 1
```

### Using Environment Variables

You can also configure the application using environment variables, which is useful in containerized environments.

```sh
export RMQ_HOST="rabbitmq.example.com"
export RMQ_USER="admin"
export RMQ_PASSWORD="super-secret-password"
export RMQ_STREAMS="telemetry,logs"
export LOG_INTERVAL="10"
export RMQ_CRC_CHECK="false"

./rmq-bench
```
