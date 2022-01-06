package grpc

import (
	"context"
	"fmt"
	"time"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcprom "github.com/grpc-ecosystem/go-grpc-prometheus"
	grpcot "github.com/opentracing-contrib/go-grpc"
	ot "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const dialTimeout = 10 * time.Second

// NewGrpcConnection dials a new traced and metered connection for gRPC client with round robin client-side LB.
func NewGrpcConnection(cfg Config, cred credentials.TransportCredentials) (*grpc.ClientConn, error) {
	tracer := ot.GlobalTracer()

	grpcprom.EnableClientHandlingTimeHistogram()

	unaryInterceptors := []grpc.UnaryClientInterceptor{
		grpcprom.UnaryClientInterceptor,
		grpcot.OpenTracingClientInterceptor(tracer),
	}

	streamInterceptors := []grpc.StreamClientInterceptor{
		grpcot.OpenTracingStreamClientInterceptor(tracer),
		grpcprom.StreamClientInterceptor,
	}

	connOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(cred),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithUnaryInterceptor(grpcmiddleware.ChainUnaryClient(unaryInterceptors...)),
		grpc.WithStreamInterceptor(grpcmiddleware.ChainStreamClient(streamInterceptors...)),
	}

	dialCtx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()

	conn, err := grpc.DialContext(dialCtx, cfg.Host, connOpts...)
	if err != nil {
		return nil, fmt.Errorf("connection dial failed: %w", err)
	}

	return conn, nil
}
