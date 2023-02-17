package main

import (
	"github.com/CrescentKohana/Zeniire/internal/config"
	grpcAPI "github.com/CrescentKohana/Zeniire/pkg/api/grpc"
	restAPI "github.com/CrescentKohana/Zeniire/pkg/api/rest"
)

// main launches the client while loading environmentals, and initializing the gRPC and REST clients.
func main() {
	config.LoadEnv()
	grpcAPI.InitGRPCClient()
	restAPI.InitRESTClient()
}
