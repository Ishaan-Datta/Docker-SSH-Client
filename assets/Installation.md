# Installation Guide

## Prerequisites

### Client:
- Go 1.19+
- gRPC and Protobuf plugins

### Server:
- Go 1.19+
- Docker (latest version recommended)
- gRPC and Protobuf plugins

### Authentication Server:
- Go 1.19+
- gRPC and Protobuf plugins

## Installation
Clone the repository and build the CLI:
```sh
git clone https://github.com/your-repo/docker-ssh-client.git
cd docker-ssh-client
go build && go install
```

Alternatively, install via Go:
```sh
go install github.com/your-repo/docker-ssh-client@latest
```

## Configuration

this should be changed to JSON config file
... add addresses for auth provider urls, IPs of other systems...

## Running the Application
### Running the Server
Run the server on the machine you want to connect to:
```sh
go run server.go
```

### Running the Client
Run the client on the machine you want to connect from:
```sh
go run client.go
```

### Running the Authentication Server
Run the authentication server on the machine you want to connect to:
```sh
go run auth-server.go
```
