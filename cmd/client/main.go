package main

import (
	"github.com/CrescentKohana/Zeniire/internal/config"
	grpcAPI "github.com/CrescentKohana/Zeniire/pkg/api/grpc"
	httpAPI "github.com/CrescentKohana/Zeniire/pkg/api/http"
)

// main launches the client while loading environmentals, and initializing the gRPC and REST clients.
func main() {
	config.LoadEnv()
	grpcAPI.InitGRPCClient()
	httpAPI.InitRESTClient()
}
