# TCP-Socket chat server

A simple chat-server using tcp-socket with TLS

## Getting Started

**1. Prerequisites:**
  
- [OpenSSl](https://www.openssl.org/)
- [Go](https://go.dev/doc/install) installed
- You also need a [tcp-chat client](https://github.com/AnhTTx13/tcp-chat) as a user interface.

**2. Install the repository:**
  
  ```sh
    go install github.com/AnhTTx13/tcp-chat-server
  ```

## Usage
  
  ```sh
    tcp-chat-server
  ```

**Flags:**

- `--host [string]`: Specify the host number.

**Note:** You must setup SSL certificates the first time you run the app.
