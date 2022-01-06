# go-grpc-client-shop
Simple example of mTLS secured gRPC client in GO that communicates with the gRPC server.
The client uses round-robin LB. Keep in mind that in case of rotating certificates the client cert has to be reloaded
with each file change.
