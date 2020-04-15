package server

import (
	context "context"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/The-Tensox/erutan/game"
	erutan "github.com/The-Tensox/erutan/protobuf"
	"github.com/The-Tensox/erutan/utils"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	keepalive "google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	status "google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
	crt                = "server1.crt"
	key                = "server1.key"
)

type Server struct {
	Host string
}

// NewServer constructor
func NewServer(host string) *Server {
	return &Server{
		Host: host,
	}
}

func (s *Server) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cert, err := credentials.NewServerTLSFromFile(crt, key)
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}

	maxConnectionAgeGrace, _ := time.ParseDuration("20s")
	t, _ := time.ParseDuration("1s")

	opts := []grpc.ServerOption{
		// The following grpc.ServerOption adds an interceptor for all unary
		// RPCs. To configure an interceptor for streaming RPCs, see:
		// https://godoc.org/google.golang.org/grpc#StreamInterceptor
		grpc.UnaryInterceptor(ensureValidToken),
		//grpc.StreamInterceptor(streamInterceptor),
		// Enable TLS for all incoming connections.
		grpc.Creds(cert),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAgeGrace: maxConnectionAgeGrace,
			Time:                  t,
		}),
	}
	srv := grpc.NewServer(opts...)
	erutan.RegisterErutanServer(srv, s)

	l, err := net.Listen("tcp", s.Host)
	if err != nil {
		return err
	}

	game.Initialize()
	go s.broadcast(ctx)
	go game.ManagerInstance.Run()

	go func() {
		err = srv.Serve(l)
		if err != nil {
			utils.ServerLogf(time.Now(), "Failed to serve gRPC server")
		}
		cancel()
	}()

	<-ctx.Done()

	close(game.ManagerInstance.Broadcast)
	utils.ServerLogf(time.Now(), "shutting down")

	srv.GracefulStop()
	return nil
}

func (s *Server) Stream(srv erutan.Erutan_StreamServer) error {
	tkn := utils.RandomString()
	// TODO: later could image a default config for clients inside a yaml or whatever
	clientSettings := erutan.Packet_UpdateParameters{
		UpdateParameters: &erutan.Packet_UpdateParametersPacket{
			Parameters: []*erutan.Packet_UpdateParametersPacket_Parameter{
				{
					Type: &erutan.Packet_UpdateParametersPacket_Parameter_TimeScale{TimeScale: 1},
				},
				{
					Type: &erutan.Packet_UpdateParametersPacket_Parameter_Debug{Debug: false},
				},
			},
		},
	}

	game.ManagerInstance.ClientsSettings.Store(tkn, clientSettings)
	go s.sendBroadcasts(srv, tkn)
	for {
		req, err := srv.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		// Distribute to the game manager to handle logic
		game.ManagerInstance.Handle(tkn, *req)
	}

	<-srv.Context().Done()
	return srv.Context().Err()
}

func (s *Server) sendBroadcasts(srv erutan.Erutan_StreamServer, tkn string) {
	stream := s.openStream(tkn)
	defer s.closeStream(tkn)

	// Notify that this client just connected
	cs, _ := game.ManagerInstance.ClientsSettings.Load(tkn)
	game.ManagerInstance.Watch.NotifyAll(utils.Event{
		Value: utils.ClientConnected{
			ClientToken: tkn,
			Settings:    cs.(erutan.Packet_UpdateParameters),
		},
	})
	for {
		select {
		case <-srv.Context().Done():
			return
		case res := <-stream:
			//if x := res.GetUpdateEntity(); x != nil {
			//	utils.DebugLogf("Sending %v", x.Components)
			//}
			if s, ok := status.FromError(srv.Send(&res)); ok {
				switch s.Code() {
				case codes.OK:
					// noop
				case codes.Unavailable, codes.Canceled, codes.DeadlineExceeded:
					utils.DebugLogf("client (%s) terminated connection", tkn)
					return
				default:
					utils.ClientLogf(time.Now(), "failed to send to client (%s): %v", tkn, s.Err())
					return
				}
			}
		}
	}
}

func (s *Server) broadcast(ctx context.Context) {
	for res := range game.ManagerInstance.Broadcast {
		game.ManagerInstance.ClientsOut.Range(func(key interface{}, stream interface{}) bool {
			if stream, ok := stream.(chan erutan.Packet); ok {
				select {
				case stream <- res:
					// noop
				default:
					utils.ServerLogf(time.Now(), "client stream full, dropping message")
				}
				return true
			}
			return false
		})
	}
}

// Initialize the communication of this client
func (s *Server) openStream(tkn string) (stream chan erutan.Packet) {
	out := make(chan erutan.Packet, 10000000) // RIP COMPUTER
	game.ManagerInstance.ClientsOut.Store(tkn, out)

	utils.ServerLogf(time.Now(), "Opened stream for client [%s]", tkn)
	// Return the channel
	return out
}

func (s *Server) closeStream(tkn string) {
	if stream, ok := game.ManagerInstance.ClientsOut.Load(tkn); ok {
		game.ManagerInstance.ClientsOut.Delete(tkn)
		if res, ok := stream.(chan erutan.Packet); ok {
			close(res)
		}
	}
	game.ManagerInstance.ClientsSettings.Delete(tkn)
	utils.ServerLogf(time.Now(), "Closed stream for client %s", tkn)
}

// valid validates the authorization.
func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	// Perform the token validation here. For the sake of this example, the code
	// here forgoes any of the usual OAuth2 token validation and instead checks
	// for a token matching an arbitrary string.
	return token == "some-secret-token"
}

// ensureValidToken ensures a valid token exists within a request's metadata. If
// the token is missing or invalid, the interceptor blocks execution of the
// handler and returns an error. Otherwise, the interceptor invokes the unary
// handler.
func ensureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}
	// The keys within metadata.MD are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	if !valid(md["authorization"]) {
		return nil, errInvalidToken
	}
	// Continue execution of handler after ensuring a valid token.
	return handler(ctx, req)
}

func streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return nil
}
