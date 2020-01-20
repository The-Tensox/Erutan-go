// Package main implements a server for Greeter service.
package main

import (
	"context"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"

	. "github.com/user/grpc_push/pkg"

	pb "github.com/user/grpc_push/protos"
	"google.golang.org/grpc"
)

const (
	port          = ":50051"
	clientAddress = "localhost:50050"
	defaultName   = "world"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// Client ...
type Client struct {
	Name string
}

var clients map[string]Client

// NewClient is called to add a new client
func NewClient(name string) {
	clients[RandomString(12)] = Client{
		name,
	}
	log.Printf("Clients: %v", clients)
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("%v trying to authentify", in.GetName())
	NewClient(in.GetName())
	return &pb.HelloReply{Message: "You are authentified " + in.GetName()}, nil
}

// SendPosition : client to server
func (s *server) SendPosition(stream pb.Greeter_SendPositionServer) error {
	for {
		position, _ := stream.Recv()
		log.Printf("Received position: (%v, %v, %v)", position.GetX(), position.GetY(), position.GetZ())
	}
}

// GetPositions ...
func (s *server) GetPositions(in *empty.Empty, stream pb.Greeter_GetPositionsServer) error {
	if err := stream.Send(&pb.Position{X: 0, Y: 0, Z: 0}); err != nil {
		return err
	}
	return nil
}

// SayBye implements helloworld.GreeterServer
func (s *server) SayBye(ctx context.Context, in *pb.ByeRequest) (*pb.ByeReply, error) {
	log.Printf("Bye %v", in.GetName())
	return &pb.ByeReply{Message: "You are disconnected " + in.GetName()}, nil
}

func main() {
	clients = make(map[string]Client)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("Server opened on port %v", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
