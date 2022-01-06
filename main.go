package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/plieskovsky/go-grpc-client-shop/internal/grpc"
	"github.com/plieskovsky/go-grpc-server-shop/proto"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
)

func main() {
	tlsCred := createClientMTLSCfg()
	conn, err := grpc.NewGrpcConnection(grpc.DefaultConfig, tlsCred)
	if err != nil {
		log.Fatalf("Failed to initialize gRPC connection with config '%+v' : %v.", grpc.DefaultConfig, err)
	}

	c := proto.NewShopServiceClient(conn)
	ctx := context.Background()

	i, err := c.Create(ctx, &proto.CreateItemRequest{
		Name:  "item name",
		Price: 55.6,
	})
	if err != nil {
		log.Fatal("Failed to create item.")
	}
	log.Printf("Created item '%+v'.", i)

	i, err = c.Get(ctx, &proto.ItemRequestId{Id: i.GetId()})
	if err != nil {
		log.Fatal("Failed to get item.")
	}
	log.Printf("Get item '%+v'.", i)

	i, err = c.Update(ctx, &proto.Item{Id: i.GetId(), Name: "update name", Price: 7894564.45})
	if err != nil {
		log.Fatal("Failed to update item.")
	}
	log.Printf("Updated item '%+v'.", i)

	i, err = c.Create(ctx, &proto.CreateItemRequest{Name: "name of second", Price: 22.22})
	if err != nil {
		log.Fatal("Failed to create second item.")
	}
	log.Printf("Created second item '%+v'.", i)

	items, err := c.GetAll(ctx, &empty.Empty{})
	if err != nil {
		log.Fatal("Failed to get all items.")
	}
	log.Printf("Got all items '%+v'.", items)

	_, err = c.Remove(ctx, &proto.ItemRequestId{Id: i.GetId()})
	if err != nil {
		log.Fatal("Failed to delete second item.")
	}
	log.Print("Deleted item.")

	items, err = c.GetAll(ctx, &empty.Empty{})
	if err != nil {
		log.Fatal("Failed to get all items.")
	}
	log.Printf("Got all items '%+v'.", items)
}

func createClientMTLSCfg() credentials.TransportCredentials {
	rootCAs, err := createCACertPool()
	if err != nil {
		log.Fatalf("Root CAs load failed: %v.", err)
	}
	cert, err := createCert()
	if err != nil {
		log.Fatalf("Client cert load failed: %v.", err)
	}
	return credentials.NewTLS(&tls.Config{
		GetClientCertificate: func(info *tls.CertificateRequestInfo) (*tls.Certificate, error) {
			return cert, nil
		},
		RootCAs: rootCAs,
	})
}

func createCert() (*tls.Certificate, error) {
	keyPair, err := tls.LoadX509KeyPair("test-certs/client-cert.pem", "test-certs/client-key.pem")
	if err != nil {
		return nil, err
	}

	return &keyPair, nil
}

func createCACertPool() (*x509.CertPool, error) {
	b, err := ioutil.ReadFile("test-certs/ca-cert.pem")
	if err != nil {
		return nil, errors.Wrap(err, "CA certificate read")
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return nil, errors.New("failed to append client CA certificate")
	}
	return cp, nil
}
