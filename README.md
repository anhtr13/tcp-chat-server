# TCP-Chat server

A simple chat server to communicate with other devices on your local network via TCP socket. Protected with transport layer security (TLS).

## Getting Started

**1. Prerequisites:**

- [OpenSSL](https://www.openssl.org/) installed.
- [Go](https://go.dev/doc/install) installed.
- You also need a [tcp-chat client](https://github.com/anhtr13/tcp-chat) as a user interface.

**2. Install the repository:**

  ```sh
    go install github.com/anhtr13/tcp-chat-server@latest
  ```

## Usage

  ```sh
    tcp-chat-server [flag]
  ```

**Flags:**

- `--port [string]`: Specify the port number (default 8080).

**Note:** You must setup SSL certificates the first time you run the app.
