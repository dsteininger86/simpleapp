package main

import (
	"context"
	"flag"
	"html/template"
	"log"
	"net/http"
	"time"

	pb "github.com/dsteininger86/simpleapp/envlookup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EnvironmentVar struct {
	Key   string
	Value string
}

var (
	backendAddr = flag.String("backend-addr", "localhost:50051", "address of the backend server")
	listenAddr  = flag.String("listen-addr", ":8080", "address to listen on")
	gRPCTimeout = flag.Duration("grpc-timeout", 5*time.Second, "timeout for gRPC calls")
	gRPCClient  pb.EnvLookupClient
	tmpl        = template.Must(template.New("index").Parse(`<!DOCTYPE html>
<html>
	<head>
		<title>Simple App</title>
	</head>
	<body>
		<h1>Simple App</h1>
		<p>The env <b>${{.Key}}</b> of the backend system is: <b>{{.Value}}</b></p>
	</body>
</html>
`))
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*backendAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("gRPC connection error: %v", err)
	}
	defer conn.Close()

	gRPCClient = pb.NewEnvLookupClient(conn)

	http.HandleFunc("/", indexHandler)
	err = http.ListenAndServe(*listenAddr, nil)
	if err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}
}

func envLookUp(ctx context.Context, env *EnvironmentVar) (bool, error) {

	ctx, cancel := context.WithTimeout(ctx, *gRPCTimeout)
	defer cancel()

	res, err := gRPCClient.GetEnv(ctx, &pb.GetEnvRequest{
		Name: env.Key,
	}, grpc.WaitForReady(true))

	if err != nil {
		return false, err
	}

	if !res.Found {
		return false, nil
	}

	env.Value = res.Value
	return true, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	env := &EnvironmentVar{
		Key: "MESSAGE",
	}

	if r.URL.Query().Get("env") != "" {
		env.Key = r.URL.Query().Get("env")
	}

	found, err := envLookUp(r.Context(), env)
	if err != nil {
		log.Printf("envLookUp error: %v", err)
		http.Error(w, "backend error", http.StatusInternalServerError)
		return
	}

	if !founded {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	err = tmpl.Execute(w, env)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
