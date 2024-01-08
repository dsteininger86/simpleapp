package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"

	pb "github.com/dsteininger86/simpleapp/envlookup"
	"google.golang.org/grpc"
)

type EnvLookupServer struct {
	pb.UnimplementedEnvLookupServer
}

func (e *EnvLookupServer) GetEnv(ctx context.Context, r *pb.GetEnvRequest) (*pb.GetEnvResponse, error) {
	envValue, found := os.LookupEnv(r.Name)
	return &pb.GetEnvResponse{
		Found: found,
		Value: envValue,
	}, nil
}

var (
	listenAddr = flag.String("listen-addr", ":50051", "address to listen on")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalf("failed to listen on %s - %v", *listenAddr, err)
	}

	serverRegistrar := grpc.NewServer()
	envLookupServer := &EnvLookupServer{}

	pb.RegisterEnvLookupServer(serverRegistrar, envLookupServer)
	err = serverRegistrar.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve %v", err)
	}
}
