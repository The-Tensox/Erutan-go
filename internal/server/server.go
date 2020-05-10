package server

import (
	"context"
	"fmt"
	"github.com/The-Tensox/Erutan-go/internal/cfg"
	"github.com/The-Tensox/Erutan-go/internal/erutan"
	"github.com/The-Tensox/Erutan-go/internal/log"
	"github.com/The-Tensox/Erutan-go/internal/obs"
	"github.com/The-Tensox/protometry"
	"github.com/golang/protobuf/ptypes"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/The-Tensox/Erutan-go/internal/utils"

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



	maxConnectionAgeGrace, _ := time.ParseDuration("20s")
	t, _ := time.ParseDuration("1s")

	opts := []grpc.ServerOption{
		// The following grpc.ServerOption adds an interceptor for all unary
		// RPCs. To configure an interceptor for streaming RPCs, see:
		// https://godoc.org/google.golang.org/grpc#StreamInterceptor
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_prometheus.StreamServerInterceptor,
			grpc_zap.StreamServerInterceptor(log.Zap),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
			grpc_zap.UnaryServerInterceptor(log.Zap),
			ensureValidToken,
		)),
		//grpc.StreamInterceptor(streamInterceptor),
		// Enable TLS for all incoming connections.
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAgeGrace: maxConnectionAgeGrace,
			Time:                  t,
		}),
	}
	if cfg.Get().Ssl {
		cert, err := credentials.NewServerTLSFromFile(cfg.Get().SslCert, cfg.Get().SslKey)
		if err != nil {
			log.Zap.Error("failed to load key pair", zap.Error(err))
		}
		opts = append(opts, grpc.Creds(cert))
	}
	srv := grpc.NewServer(opts...)
	grpc_prometheus.Register(srv)
	// Register Prometheus metrics handler.
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%s", cfg.Get().MetricsPort), nil)
	erutan.RegisterErutanServer(srv, s)
	l, err := net.Listen("tcp", s.Host)
	if err != nil {
		return err
	}

	erutan.Initialize()
	go s.broadcast(ctx)
	go erutan.ManagerInstance.Run()

	go func() {
		err = srv.Serve(l)
		if err != nil {
			log.Zap.Error("Failed to serve gRPC server")
		}
		cancel()
	}()

	<-ctx.Done()

	close(erutan.ManagerInstance.BroadcastOut)
	log.Zap.Info("shutting down")

	srv.GracefulStop()
	return nil
}

func (s *Server) Stream(srv erutan.Erutan_StreamServer) error {
	tkn := utils.RandomString()
	// TODO: later could image a default config for clients inside a yaml or whatever
	pos := cfg.Get().Logic.Player.Spawn
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
						CullingArea: protometry.NewBoxOfSize(pos[0], pos[1], pos[2], cfg.Get().Logic.Player.Culling),
					},
				},
			},
		},
	}

	erutan.ManagerInstance.ClientsSettings.Store(tkn, clientSettings)
	// Notify that this client just connected
	cs, _ := erutan.ManagerInstance.ClientsSettings.Load(tkn)

	streamOpened := make(chan bool)
	go s.clientStreamHandler(srv, tkn, streamOpened)
	<-streamOpened // Wait that the stream has been opened to notify a client connection
	erutan.ManagerInstance.NotifyAll(obs.Event{
		Value: erutan.ClientConnection{
			ClientToken: tkn,
			Settings:    cs.(erutan.Packet_UpdateParameters),
		},
	})
	for {
		req, err := srv.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		// Distribute to the erutan manager to handle logic
		erutan.ManagerInstance.ClientsIn<-*erutan.NewClientPacket(tkn, *req)
		//erutan.ManagerInstance.OnClientPacket(tkn, *req)
	}

	<-srv.Context().Done()
	return srv.Context().Err()
}

func (s *Server) clientStreamHandler(srv erutan.Erutan_StreamServer, tkn string, streamOpened chan bool) {
	stream := s.openStream(tkn)
	defer s.closeStream(tkn)
	streamOpened<-true
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
			if s, ok := status.FromError(srv.Send(&res)); ok {
				switch s.Code() {
				case codes.OK:
					// noop
				case codes.Unavailable, codes.Canceled, codes.DeadlineExceeded:
					log.Zap.Info("client terminated connection", zap.String("client", tkn))
					return
				default:
					log.Zap.Info("failed to send to client", zap.String("client", tkn), zap.Error(s.Err()))
					return
				}
			}
		}
	}
}

func (s *Server) broadcast(ctx context.Context) {
	for res := range erutan.ManagerInstance.BroadcastOut {
		//utils.DebugLogf("Broadcasting %T", res.Type)
		erutan.ManagerInstance.ClientsOut.Range(func(key, value interface{}) bool {
			if channel, ok := value.(chan erutan.Packet); ok {
				select {
				case channel <- res:
					// noop
				default:
					log.Zap.Error("client stream full, dropping message")
				}
			}
			return true
		})
	}
}

// Initialize the communication of this client
func (s *Server) openStream(tkn string) (stream chan erutan.Packet) {
	out := make(chan erutan.Packet, 10000000)
	erutan.ManagerInstance.ClientsOut.Store(tkn, out)
	log.Zap.Info("Opened stream for client", zap.String("client", tkn))
	// Return the channel
	return out
}

func (s *Server) closeStream(tkn string) {
	if inter, ok := erutan.ManagerInstance.ClientsOut.Load(tkn); ok {
		if channel, ok2 := inter.(chan erutan.Packet); ok2 {
			erutan.ManagerInstance.ClientsOut.Delete(tkn)
			close(channel)
		}
	}
	erutan.ManagerInstance.ClientsSettings.Delete(tkn)
	erutan.ManagerInstance.NotifyAll(obs.Event{ // TODO: check if that doesn't break anything since running in another goroutine
		Value: erutan.ClientDisconnection{
			ClientToken: tkn,
		},
	})
	log.Zap.Info("Closed stream for client", zap.String("client", tkn))
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
