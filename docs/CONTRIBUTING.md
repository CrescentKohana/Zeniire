# Contributing

## git

Use [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) when making commits.

## Setup

### 🚧 Building and running
#### Locally
- Initialize database PostgreSQL database
- Copy `example.env` as `.env` and change the values according to your needs
- Build server `go build ./cmd/server` and client `go build ./cmd/client`
- Run `server` (`server.exe` on Windows) and client `client` (`client.exe` on Windows)

#### Docker
- Copy `docs/docker-compose.yml` to the root of the project `./` and modify it to your needs
- Run `docker-compose up`

### 💾 Database migrations
- Automatically runs when the server is launched. Can be disabled by setting `ZNRE_DB_MIGRATIONS=false` in `.env`.

### 🔬 Testing
- Test: `go test ./... -v  -coverprofile "coverage.out"`
- Show coverage report: `go tool cover -html "coverage.out"`

### 📝 Generating docs
- Run `godoc -http=localhost:8080`
- Go to `http://localhost:8080/pkg/#thirdparty`

## Requirements
- Protocol Buffers 3.21.12+
- Go 1.20+
  - For godoc `go install golang.org/x/tools/cmd/godoc@latest`
  - For Go Protocol Buffers
    - ```
      go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
      go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
      go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest
      go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
      go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
      ```
- PostgreSQL 15+
- Docker (optional)
