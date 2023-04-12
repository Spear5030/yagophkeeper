package server

import (
	"context"
	"fmt"
	"github.com/Spear5030/yagophkeeper/internal/pb"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"time"
)

type YaGophKeeperServer struct {
	pb.UnimplementedYaGophKeeperServer
	usecase   usecase
	server    *grpc.Server
	logger    *zap.Logger
	port      string
	secretKey []byte
}

type usecase interface {
	RegisterUser(email string, password string) (token string, err error)
	LoginUser(email string, password string) (token string, err error)
	GetLastSyncTime(email string) (lastSync time.Time, err error)
	SetData(email string, data []byte) (err error)
	GetData(email string) (data []byte, err error)
}

func New(usecase usecase, logger *zap.Logger, port string) *YaGophKeeperServer {
	s := &YaGophKeeperServer{
		usecase: usecase,
		logger:  logger,
		port:    port,
	}
	s.server = grpc.NewServer(grpc.UnaryInterceptor(s.AuthInterceptor))
	reflection.Register(s.server) // for postman
	pb.RegisterYaGophKeeperServer(s.server, s)
	s.secretKey = []byte("secret") //todo config
	return s
}

// Start слушает определенный порт и запускает в горутине grpc сервер
func (s *YaGophKeeperServer) Start() error {
	l, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		log.Fatal("error with listen gRPC:", err)
	}
	return s.server.Serve(l)
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

func (s *YaGophKeeperServer) CheckSync(ctx context.Context, req *pb.CheckSyncRequest) (*pb.SyncResponse, error) {
	var resp = &pb.SyncResponse{}
	fmt.Println(req)
	email := getEmailFromContext(ctx)
	s.logger.Debug(email)
	lastSync, err := s.usecase.GetLastSyncTime(email)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp.LastSync = timestamppb.New(lastSync)
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

func (s *YaGophKeeperServer) SetData(ctx context.Context, secrets *pb.Secrets) (*pb.SyncResponse, error) {
	var resp = &pb.SyncResponse{}
	err := s.usecase.SetData(getEmailFromContext(ctx), secrets.Data)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp.LastSync = timestamppb.New(time.Now())
	return resp, err
}

func (s *YaGophKeeperServer) GetData(ctx context.Context, empty *emptypb.Empty) (*pb.Secrets, error) {
	var resp = &pb.Secrets{}
	var err error
	resp.Data, err = s.usecase.GetData(getEmailFromContext(ctx))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return resp, nil
}

func (s *YaGophKeeperServer) AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	switch info.FullMethod {
	case "/yagophkeeper.YaGophKeeper/RegisterUser":
		return handler(ctx, req)
	case "/yagophkeeper.YaGophKeeper/LoginUser":
		return handler(ctx, req)
	}
	var token *jwt.Token
	var err error
	s.logger.Debug("auth interceptor")
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("Bearer")
		if len(values) > 0 {
			s.logger.Debug(values[0])
			token, err = jwt.Parse(values[0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return s.secretKey, nil
			})
			if err != nil {
				s.logger.Debug(err.Error())
				return nil, status.Error(codes.Unauthenticated, err.Error())
			}
		}
	}
	if claims, ok := token.Claims.(jwt.MapClaims); token.Valid && ok {
		email := claims["email"].(string)
		s.logger.Debug("user email from jwt", zap.String("email", email))
		ctx = metadata.AppendToOutgoingContext(ctx, "email", email) //todo check merged keys
		return handler(ctx, req)
	} else {
		return nil, status.Error(codes.Unauthenticated, "wrong token") //todo check exp
	}
}

func getEmailFromContext(ctx context.Context) (email string) {
	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		values := md.Get("email")
		if len(values) > 0 {
			email = values[0]
		}
	}
	return
}
