package auth

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/CrescentKohana/Zeniire/internal/config"
	"google.golang.org/grpc/credentials"
	"os"
)

func CreateClientCertPool() (*x509.CertPool, error) {
	pemServerCA, err := os.ReadFile(config.Options.GRPC.CertsPath + "/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	return certPool, nil
}

// LoadServerTLSCredentials loads server's certificate and private key, and then creates the credentials and returns them
func LoadServerTLSCredentials() (credentials.TransportCredentials, error) {
	certPool, err := CreateClientCertPool()
	if err != nil {
		return nil, err
	}

	serverCert, err := tls.LoadX509KeyPair(
		config.Options.GRPC.CertsPath+"/server-cert.pem",
		config.Options.GRPC.CertsPath+"/server-key.pem",
	)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(tlsConfig), nil
}

// LoadClientTLSCredentials loads the certificate of the CA who signed server's certificate,
// and then creates the credentials and returns them
func LoadClientTLSCredentials() (credentials.TransportCredentials, error) {
	certPool, err := CreateClientCertPool()
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(tlsConfig), nil
}
