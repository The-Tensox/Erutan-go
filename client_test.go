package main

import (
	context "context"
	"flag"
	"io"
	"log"
	"os"
	"testing"

	erutan "github.com/user/erutan_two/protos/realtime"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func initt() {
	os.Args = append(os.Args, "-s")
	os.Args = append(os.Args, "-v")
	os.Args = append(os.Args, "-p \"\"")
	os.Args = append(os.Args, "-h \"0.0.0.0:50051\"")
	flag.Parse()
}

func google() (*grpc.ClientConn, error) {
	auth, _ := oauth.NewApplicationDefault(context.Background(), "")
	return grpc.Dial(
		"greeter.googleapis.com", grpc.WithPerRPCCredentials(auth),
	)
}

func ssl() (*grpc.ClientConn, error) {
	creds, _ := credentials.NewClientTLSFromFile("server1.crt", "")
	return grpc.Dial(
		"127.0.0.1:50051", grpc.WithTransportCredentials(creds),
	)
}

func TestClient(t *testing.T) {
	//initt()
	channel, _ := ssl()
	client := erutan.NewErutanClient(channel)
	stream, _ := client.Stream(context.Background())

	for {
		_, err := stream.Recv()
		if err == io.EOF {
			// read done.
			return
		}
		if err != nil {
			log.Fatalf("Failed to receive a note : %v", err)
		}
	}
}
