package test

import (
	context "context"
	"io"
	"log"
	"testing"

	main "github.com/user/erutan"
	erutan "github.com/user/erutan/protos/realtime"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func setFlags() {
	Config.DebugMode = true
	Config.Host = "0.0.0.0:50051"
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
	setFlags()
	go main.RunMain()

	// TODO: panic: runtime error: invalid memory address or nil pointer dereference [recovered]
	//		 panic: runtime error: invalid memory address or nil pointer dereference
	//		 [signal SIGSEGV: segmentation violation code=0x1 addr=0x30 pc=0x851a9b]
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
