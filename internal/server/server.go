package server

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/Spear5030/yagophkeeper/internal/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

type YaGophKeeperServer struct {
	pb.UnimplementedYaGophKeeperServer
	usecase   usecase
	Server    *grpc.Server
	logger    *zap.Logger
	port      string
	secretKey string
}

type usecase interface {
	RegisterUser(email string, password string) (token string, err error)
	LoginUser(email string, password string) (token string, err error)
}

func New(usecase usecase, logger *zap.Logger, port string) *YaGophKeeperServer {
	s := &YaGophKeeperServer{
		usecase: usecase,
		logger:  logger,
		port:    port,
	}
	s.Server = grpc.NewServer(grpc.UnaryInterceptor(s.AuthInterceptor))
	reflection.Register(s.Server) // for postman
	pb.RegisterYaGophKeeperServer(s.Server, s)
	return s
}

// Start слушает определенный порт и запускает в горутине grpc сервер
func (s *YaGophKeeperServer) Start() error {
	l, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		log.Fatal("error with listen gRPC:", err)
	}
	return s.Server.Serve(l)
}

func (s *YaGophKeeperServer) RegisterUser(ctx context.Context, user *pb.User) (*pb.AuthResponse, error) {
	var resp = &pb.AuthResponse{}
	var err error
	resp.Token, err = s.usecase.RegisterUser(user.Email, user.Password)
	if err != nil {
		s.logger.Debug("RegisterUser error", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return resp, err
}

func (s *YaGophKeeperServer) LoginUser(ctx context.Context, user *pb.User) (*pb.AuthResponse, error) {
	var resp = &pb.AuthResponse{}
	var err error
	resp.Token, err = s.usecase.LoginUser(user.Email, user.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return resp, err
}

func (s *YaGophKeeperServer) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	switch info.FullMethod {
	case "/yagophkeeper.YaGophKeeper/RegisterUser":
		return handler(ctx, req)
	case "/yagophkeeper.YaGophKeeper/LoginUser":
		return handler(ctx, req)
	}
	var id, token []byte
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("token")
		if len(values) > 0 {
			token, _ = hex.DecodeString(values[0])
		}
		values = md.Get("id")
		if len(values) > 0 {
			id = []byte(values[0])
		}
	}
	if len(token) == 0 || len(id) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}
	h := hmac.New(sha256.New, []byte(s.secretKey))
	h.Write(id)

	if !hmac.Equal(h.Sum(nil), token) {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}
	return handler(ctx, req)
}
