package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

// OptionsModel is the model for the global configuration.
type OptionsModel struct {
	RESTAddr string
	GRPC     GRPCOptions
	DB       DBOptions
}

// GRPCOptions is the model for gRPC configuration.
type GRPCOptions struct {
	Host      string
	Port      string
	TLS       bool
	CertsPath string
}

// Options stores the global configuration for the application.
var Options *OptionsModel

// LoadEnv loads the environment variables into Options.
func LoadEnv() {
	var err = godotenv.Load()
	if err != nil {
		log.Debug("No .env file found")
	}

	Options = &OptionsModel{
		RESTAddr: rESTAddr(),
		GRPC: GRPCOptions{
			Host:      gRPCHost(),
			Port:      gRPCPort(),
			TLS:       gRPCTLS(),
			CertsPath: certsPath(),
		},
		DB: DBOptions{
			Migrations: dbMigrationsEnabled(),
			Address:    "postgres://" + dbUser() + ":" + dbPassword() + "@" + dbHost() + ":" + dbPort() + "/" + dbName(),
		},
	}
}

// gRPCHost loads and parses the gRPC host env.
func gRPCHost() string {
	value := os.Getenv("ZNRE_GRPC_HOST")
	if value == "" {
		return "localhost"
	}
	return value
}

// gRPCPort loads and parses the gRPC port env.
func gRPCPort() string {
	value := os.Getenv("ZNRE_GRPC_PORT")
	if value == "" {
		return "50051"
	}
	return value
}

// gRPCTLS checks if TLS is enabled through the env.
func gRPCTLS() bool {
	return os.Getenv("ZNRE_GRPC_TLS") == "true"
}

// certsPath loads and returns the path for the certificates
func certsPath() string {
	return os.Getenv("ZNRE_CERTS_PATH")
}

// rESTAddr loads and parses both the REST host and port envs.
func rESTAddr() string {
	host := os.Getenv("ZNRE_REST_HOST")
	port := os.Getenv("ZNRE_REST_PORT")
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "3333"
	}
	return host + ":" + port
}
