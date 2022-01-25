# guardrails-assignment

## Prerequisites
Before you continue, ensure you meet the following requirements:

### Dependencies
```sh
- Docker
- Go 1.17
- PostgreSQL 12
```

## Create instance database
```sh
make init-db
```

## Start server
Default server always start in `port: 10080`. Please change configure in `/config/default.go` if you want to start with new port
```sh
go run cmd/main.go server
```

## Start scanner job
```sh
go run cmd/main.go scanner
```

## Installation buf
The following instructions assume you are using [Buf Docs Installation](https://docs.buf.build/installation)
```sh
brew tap bufbuild/buf
brew install buf
```


## Installation grpc

The following instructions assume you are using
[Go Modules](https://github.com/golang/go/wiki/Modules) for dependency
management. Use a
[tool dependency](https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module)
to track the versions of the following executable packages:

```go
// +build tools

package tools

import (
    _ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
    _ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
    _ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
    _ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
```

Run `go mod tidy` to resolve the versions.

Install by running
```sh
go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

This will place four binaries in your `$GOBIN`;

- `protoc-gen-grpc-gateway`
- `protoc-gen-openapiv2`
- `protoc-gen-go`
- `protoc-gen-go-grpc`

Make sure that your `$GOBIN` is in your `$PATH`.

Alternatively, see the section on remotely managed plugin versions below.

## Generate code

```sh
buf generate
```