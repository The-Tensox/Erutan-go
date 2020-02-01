package main

import (
	context "context"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	erutan "github.com/user/erutan_two/protos/realtime"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	keepalive "google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	status "google.golang.org/grpc/status"
)

const tokenHeader = "x-token"

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
	crt                = "server1.crt"
	key                = "server1.key"
	// Broadcast is a global channel to send packets to clients
	Broadcast chan erutan.Packet = make(chan erutan.Packet, 1000)

	// TickRate defines the Server's tick rate, the lower the faster
	TickRate float64 = 35
)

type Server struct {
	Host, Password string

	ClientNames   map[string]string
	ClientStreams map[string]chan erutan.Packet

	namesMtx, streamsMtx sync.RWMutex
}

// NewServer constructor
func NewServer(host, pass string) *Server {
	return &Server{
		Host:     host,
		Password: pass,

		ClientNames:   make(map[string]string),
		ClientStreams: make(map[string]chan erutan.Packet),
	}
}

func (s *Server) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ServerLogf(time.Now(),
		"starting on %s with password %q",
		s.Host, s.Password)

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

	go s.broadcast(ctx)

	go RunGame()

	go func() {
		srv.Serve(l)
		cancel()
	}()

	<-ctx.Done()

	close(Broadcast)
	ServerLogf(time.Now(), "shutting down")

	srv.GracefulStop()
	return nil
}

func (s *Server) Stream(srv erutan.Erutan_StreamServer) error {
	tkn := RandomString()
	go s.sendBroadcasts(srv, tkn)
	for {
		req, err := srv.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		DebugLogf("client send: %v", req)
		// Use a type switch to determine which oneof was set.
		/*
			switch t := req.Type.(type) {
			case *erutan.ClientToServer_UpdatePosition:
				s.updatePosition(req.GetUpdatePosition())
			default:
				DebugLogf("client send: unimplemented packet handler: %v", t)
			}
		*/
	}

	<-srv.Context().Done()
	return srv.Context().Err()
}

func (s *Server) sendBroadcasts(srv erutan.Erutan_StreamServer, tkn string) {
	stream := s.openStream(tkn)
	defer s.closeStream(tkn)

	// Send world state (only) to the client that connected
	for _, packet := range WorldState() {
		srv.Send(packet)
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
					DebugLogf("client (%s) terminated connection", tkn)
					return
				default:
					ClientLogf(time.Now(), "failed to send to client (%s): %v", tkn, s.Err())
					return
				}
			}
		}
	}
}

func (s *Server) broadcast(ctx context.Context) {
	for res := range Broadcast {
		s.streamsMtx.RLock()
		for _, stream := range s.ClientStreams {
			select {
			case stream <- res:
				// noop
			default:
				ServerLogf(time.Now(), "client stream full, dropping message")
			}
		}
		s.streamsMtx.RUnlock()
	}
}

func (s *Server) openStream(tkn string) (stream chan erutan.Packet) {
	stream = make(chan erutan.Packet, 100)

	s.streamsMtx.Lock()
	s.ClientStreams[tkn] = stream
	s.streamsMtx.Unlock()

	DebugLogf("opened stream for client %s", tkn)

	return
}

func (s *Server) closeStream(tkn string) {
	s.streamsMtx.Lock()

	if stream, ok := s.ClientStreams[tkn]; ok {
		delete(s.ClientStreams, tkn)
		close(stream)
	}

	DebugLogf("closed stream for client %s", tkn)

	s.streamsMtx.Unlock()
}

func (s *Server) getName(tkn string) (name string, ok bool) {
	s.namesMtx.RLock()
	name, ok = s.ClientNames[tkn]
	s.namesMtx.RUnlock()
	return
}

func (s *Server) setName(tkn string, name string) {
	s.namesMtx.Lock()
	s.ClientNames[tkn] = name
	s.namesMtx.Unlock()
}

func (s *Server) delName(tkn string) (name string, ok bool) {
	name, ok = s.getName(tkn)

	if ok {
		s.namesMtx.Lock()
		delete(s.ClientNames, tkn)
		s.namesMtx.Unlock()
	}

	return
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

	/*s, _ := srv.(*Server)

	s.streamsMtx.Lock()

	DebugLogf("streamInterceptor ", s.ClientStreams, ss, info, handler)
	s.streamsMtx.Unlock()
	*/
	return nil
}
