package grpc

import (
	"time"
)

const defaultTimeout = 500 * time.Millisecond

var DefaultConfig = Config{
	// https://rafaeleyng.github.io/grpc-load-balancing-with-grpc-go for DNS LB
	Host:    "dns:///localhost:8443",
	Timeout: defaultTimeout,
}

// Config is the gRPC client configuration.
type Config struct {
	Host    string
	Timeout time.Duration
}
