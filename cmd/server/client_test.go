package main

import (
	"fmt"
	"github.com/The-Tensox/Erutan-go/internal/cfg"
	"github.com/The-Tensox/Erutan-go/internal/erutan"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"testing"
	"time"
)

//func google() (*grpc.ClientConn, error) {
//	auth, _ := oauth.NewApplicationDefault(context.Background(), "")
//	return grpc.Dial(
//		"greeter.googleapis.com", grpc.WithPerRPCCredentials(auth),
//	)
//}


func connect(t *testing.T) *grpc.ClientConn {
	go RunMain()
	tls := true
	var opts []grpc.DialOption
	if tls {
		creds, err := credentials.NewClientTLSFromFile(cfg.Get().SslCert, "")
		if err != nil {
			t.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithBlock())
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.Get().Server.Host, cfg.Get().Server.Port), opts...)
	if err != nil {
		t.Fatalf("fail to dial: %v", err)
	}
	return conn
}

func TestClient(t *testing.T) {
	conn := connect(t)
	defer conn.Close()
	client := erutan.NewErutanClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c, err := client.Stream(ctx)
	if err != nil {
		t.Fatalf("Couldn't open stream : %v", err)
	}
	for {
		_, err := c.Recv()
		if err == io.EOF {
			// read done.
			return
		}
		if err != nil {
			t.Fatalf("Failed to receive : %v", err)
		}
		t.SkipNow()
	}
}
