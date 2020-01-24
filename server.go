package main

import (
	context "context"
	"crypto/rand"
	fmt "fmt"
	"io"
	"net"
	"sync"
	"time"

	erutan "github.com/user/erutan_two/protos/realtime"

	"github.com/golang/protobuf/ptypes"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	status "google.golang.org/grpc/status"
)

const tokenHeader = "x-token"

type server struct {
	Host, Password string

	Broadcast chan erutan.StreamResponse

	ClientNames   map[string]string
	ClientStreams map[string]chan erutan.StreamResponse

	namesMtx, streamsMtx sync.RWMutex
}

// Server constructor
func Server(host, pass string) *server {
	return &server{
		Host:     host,
		Password: pass,

		Broadcast: make(chan erutan.StreamResponse, 1000),

		ClientNames:   make(map[string]string),
		ClientStreams: make(map[string]chan erutan.StreamResponse),
	}
}

func (s *server) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ServerLogf(time.Now(),
		"starting on %s with password %q",
		s.Host, s.Password)

	srv := grpc.NewServer()
	erutan.RegisterErutanServer(srv, s)

	l, err := net.Listen("tcp", s.Host)
	if err != nil {
		return err
	}

	go s.broadcast(ctx)

	go func() {
		srv.Serve(l)
		cancel()
	}()

	<-ctx.Done()

	s.Broadcast <- erutan.StreamResponse{
		Timestamp: ptypes.TimestampNow(),
		Event: &erutan.StreamResponse_ServerShutdown{
			&erutan.StreamResponse_Shutdown{}}}

	close(s.Broadcast)
	ServerLogf(time.Now(), "shutting down")

	srv.GracefulStop()
	return nil
}

func (s *server) Login(ctx context.Context, req *erutan.LoginRequest) (*erutan.LoginResponse, error) {
	switch {
	case req.Password != s.Password:
		return nil, status.Error(codes.Unauthenticated, "password is incorrect")
	case req.Name == "":
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}

	tkn := s.genToken()
	s.setName(tkn, req.Name)

	ServerLogf(time.Now(), "%s (%s) has logged in", tkn, req.Name)

	s.Broadcast <- erutan.StreamResponse{
		Timestamp: ptypes.TimestampNow(),
		Event: &erutan.StreamResponse_ClientLogin{&erutan.StreamResponse_Login{
			Name: req.Name,
		}},
	}

	return &erutan.LoginResponse{Token: tkn}, nil
}

func (s *server) Logout(ctx context.Context, req *erutan.LogoutRequest) (*erutan.LogoutResponse, error) {
	name, ok := s.delName(req.Token)
	if !ok {
		return nil, status.Error(codes.NotFound, "token not found")
	}

	ServerLogf(time.Now(), "%s (%s) has logged out", req.Token, name)

	s.Broadcast <- erutan.StreamResponse{
		Timestamp: ptypes.TimestampNow(),
		Event: &erutan.StreamResponse_ClientLogout{&erutan.StreamResponse_Logout{
			Name: name,
		}},
	}

	return new(erutan.LogoutResponse), nil
}

func (s *server) Stream(srv erutan.Erutan_StreamServer) error {
	tkn, ok := s.extractToken(srv.Context())
	if !ok {
		return status.Error(codes.Unauthenticated, "missing token header")
	}

	name, ok := s.getName(tkn)
	if !ok {
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	go s.sendBroadcasts(srv, tkn)

	for {
		req, err := srv.Recv()
		DebugLogf("client send %v", req)

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		s.Broadcast <- erutan.StreamResponse{
			Timestamp: ptypes.TimestampNow(),
			Event: &erutan.StreamResponse_ClientMessage{&erutan.StreamResponse_Message{
				Name:    name,
				Message: req.Message,
			}},
		}
	}

	<-srv.Context().Done()
	return srv.Context().Err()
}

func (s *server) sendBroadcasts(srv erutan.Erutan_StreamServer, tkn string) {
	stream := s.openStream(tkn)
	defer s.closeStream(tkn)

	go func() {
		for range time.Tick(2000 * time.Millisecond) {
			srv.Send(&erutan.StreamResponse{
				Timestamp: ptypes.TimestampNow(),
				Event: &erutan.StreamResponse_ClientMessage{&erutan.StreamResponse_Message{
					Name:    "Server",
					Message: fmt.Sprintf("Welcome %s", s.ClientNames[tkn]),
				}},
			})
		}
	}()

	for {
		select {
		case <-srv.Context().Done():
			return
		case res := <-stream:
			if s, ok := status.FromError(srv.Send(&res)); ok {
				switch s.Code() {
				case codes.OK:
					// noop
					DebugLogf("res %v", res)

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

func (s *server) broadcast(ctx context.Context) {
	for res := range s.Broadcast {
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

func (s *server) openStream(tkn string) (stream chan erutan.StreamResponse) {
	stream = make(chan erutan.StreamResponse, 100)

	s.streamsMtx.Lock()
	s.ClientStreams[tkn] = stream
	s.streamsMtx.Unlock()

	DebugLogf("opened stream for client %s", tkn)

	return
}

func (s *server) closeStream(tkn string) {
	s.streamsMtx.Lock()

	if stream, ok := s.ClientStreams[tkn]; ok {
		delete(s.ClientStreams, tkn)
		close(stream)
	}

	DebugLogf("closed stream for client %s", tkn)

	s.streamsMtx.Unlock()
}

func (s *server) genToken() string {
	tkn := make([]byte, 4)
	rand.Read(tkn)
	return fmt.Sprintf("%x", tkn)
}

func (s *server) getName(tkn string) (name string, ok bool) {
	s.namesMtx.RLock()
	name, ok = s.ClientNames[tkn]
	s.namesMtx.RUnlock()
	return
}

func (s *server) setName(tkn string, name string) {
	s.namesMtx.Lock()
	s.ClientNames[tkn] = name
	s.namesMtx.Unlock()
}

func (s *server) delName(tkn string) (name string, ok bool) {
	name, ok = s.getName(tkn)

	if ok {
		s.namesMtx.Lock()
		delete(s.ClientNames, tkn)
		s.namesMtx.Unlock()
	}

	return
}

func (s *server) extractToken(ctx context.Context) (tkn string, ok bool) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok || len(md[tokenHeader]) == 0 {
		return "", false
	}

	return md[tokenHeader][0], true
}
