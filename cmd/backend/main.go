package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	pb "github.com/dsteininger86/simpleapp/envlookup"
	"google.golang.org/grpc"
)

type EnvLookupServer struct {
	pb.UnimplementedEnvLookupServer
}

func lookupEnv(envKey string) (string, bool) {

	envVars := os.Environ()

	for _, envVar := range envVars {
		if strings.HasPrefix(fmt.Sprintf("%s=", envVar), envKey) {
			return strings.Split(envVar, "=")[1], true
		}
	}
	return "", false
}

func (e *EnvLookupServer) GetEnv(ctx context.Context, r *pb.GetEnvRequest) (*pb.GetEnvResponse, error) {
	envValue, found := lookupEnv(r.Name)
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
