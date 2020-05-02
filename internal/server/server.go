package server

import (
	"context"
	"fmt"
	"github.com/The-Tensox/Erutan-go/internal/cfg"
	"github.com/The-Tensox/Erutan-go/internal/obs"
	"github.com/The-Tensox/protometry"
	"github.com/golang/protobuf/ptypes"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/The-Tensox/Erutan-go/internal/game"
	"github.com/The-Tensox/Erutan-go/internal/utils"
	erutan "github.com/The-Tensox/Erutan-go/protobuf"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
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

	cert, err := credentials.NewServerTLSFromFile(cfg.Global.SslCert, cfg.Global.SslKey)
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
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		//grpc.StreamInterceptor(streamInterceptor),
		// Enable TLS for all incoming connections.
		grpc.Creds(cert),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAgeGrace: maxConnectionAgeGrace,
			Time:                  t,
		}),
	}
	srv := grpc.NewServer(opts...)
	grpc_prometheus.Register(srv)
	// Register Prometheus metrics handler.
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%s", cfg.Global.MetricsPort), nil)
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
	pos := cfg.Global.Logic.Player.Spawn
	clientSettings := erutan.Packet_UpdateParameters{
		UpdateParameters: &erutan.Packet_UpdateParametersPacket{
			Parameters: []*erutan.Packet_UpdateParametersPacket_Parameter{
				{
					Type: &erutan.Packet_UpdateParametersPacket_Parameter_TimeScale{TimeScale: 1},
				},
				{
					Type: &erutan.Packet_UpdateParametersPacket_Parameter_Debug{Debug: false},
				},
				{
					Type: &erutan.Packet_UpdateParametersPacket_Parameter_CullingArea{
						CullingArea: protometry.NewBoxOfSize(pos.X, pos.Y, pos.Z, cfg.Global.Logic.Player.Culling),
					},
				},
			},
		},
	}

	game.ManagerInstance.ClientsSettings.Store(tkn, clientSettings)
	// Notify that this client just connected
	cs, _ := game.ManagerInstance.ClientsSettings.Load(tkn)
	game.ManagerInstance.Watch.NotifyAll(obs.Event{
		Value: obs.ClientConnection{
			ClientToken: tkn,
			Settings:    cs.(erutan.Packet_UpdateParameters),
		},
	})
	go s.sendBroadcasts(srv, tkn)
	for {
		req, err := srv.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		// Distribute to the game manager to handle logic
		game.ManagerInstance.OnClientPacket(tkn, *req)
	}

	<-srv.Context().Done()
	return srv.Context().Err()
}

func (s *Server) sendBroadcasts(srv erutan.Erutan_StreamServer, tkn string) {
	stream := s.openStream(tkn)
	defer s.closeStream(tkn)

	// Tell to this client what is its token
	stream <- erutan.Packet{
		Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
		Type:     &erutan.Packet_Authentication{
			Authentication: &erutan.Packet_AuthenticationPacket{
				ClientToken: tkn,
			},
		},
	}

	for {
		select {
		case <-srv.Context().Done():
			return
		case res := <-stream:
			//utils.DebugLogf("Send")
			//if x := res.GetUpdateObject(); x != nil {
			//	utils.DebugLogf("Sending %v", x.Components)
			//}
			//utils.DebugLogf("Send")
			//if x := res.GetCreatePlayer(); x != nil {
			//	utils.DebugLogf("Sending %v", x.ObjectId)
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
		for _, stream := range game.ManagerInstance.ClientsOut {
			//if x := res.GetCreatePlayer(); x != nil {
			//	utils.DebugLogf("Sending %v", x.ObjectId)
			//}
			select {
			case stream <- res:
				// noop
			default:
				utils.ServerLogf(time.Now(), "client stream full, dropping message")
			}
		}
	}
}

// Initialize the communication of this client
func (s *Server) openStream(tkn string) (stream chan erutan.Packet) {
	out := make(chan erutan.Packet, 10000000) // RIP COMPUTER
	game.ManagerInstance.ClientsMu.Lock()
	defer game.ManagerInstance.ClientsMu.Unlock()

	game.ManagerInstance.ClientsOut[tkn] = out

	utils.ServerLogf(time.Now(), "Opened stream for client [%s]", tkn)
	// Return the channel
	return out
}

func (s *Server) closeStream(tkn string) {
	game.ManagerInstance.ClientsMu.Lock()
	defer game.ManagerInstance.ClientsMu.Unlock()
	if stream, ok := game.ManagerInstance.ClientsOut[tkn]; ok {
		delete(game.ManagerInstance.ClientsOut, tkn)
		close(stream)
	}
	game.ManagerInstance.ClientsSettings.Delete(tkn)
	game.ManagerInstance.Watch.NotifyAll(obs.Event{ // TODO: check if that doesn't break anything since running in another goroutine
		Value: obs.ClientDisconnection{
			ClientToken: tkn,
		},
	})
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
	_, _ = grpc_prometheus.UnaryServerInterceptor(ctx, req, info, handler)
	// The keys within metadata.MD are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	if !valid(md["authorization"]) {
		return nil, errInvalidToken
	}
	// Continue execution of handler after ensuring a valid token.
	return handler(ctx, req)
}
