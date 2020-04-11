package main

import (
	erutan "github.com/The-Tensox/erutan/protobuf"
	"github.com/The-Tensox/erutan/utils"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"log"
	"testing"
	"time"
)

func setFlags() {
	utils.Config.DebugMode = true
	utils.Config.Host = "0.0.0.0:50051"
}

//func google() (*grpc.ClientConn, error) {
//	auth, _ := oauth.NewApplicationDefault(context.Background(), "")
//	return grpc.Dial(
//		"greeter.googleapis.com", grpc.WithPerRPCCredentials(auth),
//	)
//}

func ssl() (*grpc.ClientConn, error) {
	creds, _ := credentials.NewClientTLSFromFile("server1.crt", "")
	return grpc.Dial(
		"127.0.0.1:50051", grpc.WithTransportCredentials(creds),
	)
}

func TestClient(t *testing.T) {
	setFlags()
	go RunMain()
	time.Sleep(10 * time.Second)     // sleep for 10 seconds
	tls := true
	crtFile := "server1.crt"
	serverAddr := "127.0.0.1:50051"
	var opts []grpc.DialOption
	if tls {
		creds, err := credentials.NewClientTLSFromFile(crtFile, "")
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := erutan.NewErutanClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c, err := client.Stream(ctx)
	if err != nil {
		log.Fatalf("Couldn't open stream : %v", err)
	}

	for {
		d, err := c.Recv()
		if err == io.EOF {
			// read done.
			return
		}
		if err != nil {
			log.Fatalf("Failed to receive : %v", err)
		}
		log.Printf("recv %v", d)
	}
}
