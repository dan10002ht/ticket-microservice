// Package grpctls provides gRPC TLS credentials toggled by the GRPC_TLS_ENABLED
// environment variable. When disabled (default), servers and clients use plaintext.
package grpctls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultCert = "/certs/server.crt"
	defaultKey  = "/certs/server.key"
	defaultCA   = "/certs/ca.crt"
)

// IsEnabled returns true when GRPC_TLS_ENABLED is "true".
func IsEnabled() bool {
	return os.Getenv("GRPC_TLS_ENABLED") == "true"
}

func certPath() string { return env("GRPC_TLS_CERT", defaultCert) }
func keyPath() string  { return env("GRPC_TLS_KEY", defaultKey) }
func caPath() string   { return env("GRPC_TLS_CA", defaultCA) }

// ServerOption returns a grpc.ServerOption that enables TLS if GRPC_TLS_ENABLED
// is "true", otherwise returns an empty (no-op) option.
func ServerOption() grpc.ServerOption {
	if !IsEnabled() {
		return grpc.EmptyServerOption{}
	}

	cert, err := tls.LoadX509KeyPair(certPath(), keyPath())
	if err != nil {
		panic(fmt.Sprintf("grpctls: failed to load server cert/key: %v", err))
	}

	ca, err := os.ReadFile(caPath())
	if err != nil {
		panic(fmt.Sprintf("grpctls: failed to read CA cert: %v", err))
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(ca) {
		panic("grpctls: failed to parse CA cert")
	}

	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    pool,
		ClientAuth:   tls.RequireAndVerifyClientCert, // mTLS
		MinVersion:   tls.VersionTLS12,
	}

	return grpc.Creds(credentials.NewTLS(tlsCfg))
}

// DialOption returns a grpc.DialOption for client connections. Uses TLS with
// the CA cert when enabled, otherwise returns insecure credentials.
func DialOption() grpc.DialOption {
	if !IsEnabled() {
		return grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	cert, err := tls.LoadX509KeyPair(certPath(), keyPath())
	if err != nil {
		panic(fmt.Sprintf("grpctls: failed to load client cert/key: %v", err))
	}

	ca, err := os.ReadFile(caPath())
	if err != nil {
		panic(fmt.Sprintf("grpctls: failed to read CA cert: %v", err))
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(ca) {
		panic("grpctls: failed to parse CA cert")
	}

	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
		MinVersion:   tls.VersionTLS12,
	}

	return grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg))
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
