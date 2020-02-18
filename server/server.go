package server

import (
	context "context"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/user/erutan/game"
	erutan "github.com/user/erutan/protos/realtime"
	"github.com/user/erutan/utils"

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

	utils.ServerLogf(time.Now(), "Server started with config: %+v", utils.Config)

	game.Initialize()
	go s.broadcast(ctx)
	go game.ManagerInstance.Run()

	go func() {
		srv.Serve(l)
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
	go s.sendBroadcasts(srv, tkn)
	for {
		req, err := srv.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		// Distribute to the game manager to handle logic
		game.ManagerInstance.Receive <- req
	}

	<-srv.Context().Done()
	return srv.Context().Err()
}

func (s *Server) sendBroadcasts(srv erutan.Erutan_StreamServer, tkn string) {
	stream := s.openStream(tkn)
	defer s.closeStream(tkn)

	// Send world state (only) to the client that connected
	game.ManagerInstance.SyncNewClient(tkn)
	for {
		select {
		case <-srv.Context().Done():
			return
		case res := <-stream:
			//utils.DebugLogf("Sending %v", res)
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
		game.ManagerInstance.ClientStreams.Range(func(key interface{}, stream interface{}) bool {
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

func (s *Server) openStream(tkn string) (stream chan erutan.Packet) {
	stream = make(chan erutan.Packet, 10000000) // RIP COMPUTER

	game.ManagerInstance.ClientStreams.Store(tkn, stream)

	utils.DebugLogf("opened stream for client %s", tkn)

	// Return the channel
	return
}

func (s *Server) closeStream(tkn string) {
	if stream, ok := game.ManagerInstance.ClientStreams.Load(tkn); ok {
		game.ManagerInstance.ClientStreams.Delete(tkn)
		if res, ok := stream.(chan erutan.Packet); ok {
			close(res)
		}
	}

	utils.DebugLogf("closed stream for client %s", tkn)
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
