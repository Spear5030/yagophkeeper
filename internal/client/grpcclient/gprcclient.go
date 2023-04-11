package grpcclient

import (
	"context"
	"github.com/Spear5030/yagophkeeper/internal/domain"
	"github.com/Spear5030/yagophkeeper/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"time"
)

type Client struct {
	conn       *grpc.ClientConn
	yagkclient pb.YaGophKeeperClient
	token      string
}

func New(addr string, token string) *Client {
	var c Client
	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(c.AuthInterceptor()),
	)
	if err != nil {
		log.Fatal(err)
	}
	c.conn = conn
	c.yagkclient = pb.NewYaGophKeeperClient(conn)
	c.token = token
	return &c
}

func (c *Client) RegisterUser(user domain.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := c.yagkclient.RegisterUser(ctx, &pb.User{Email: user.Email, Password: user.Password})
	if err != nil {
		return "", err
	}
	c.token = resp.Token
	return resp.Token, nil
}

func (c *Client) LoginUser(user domain.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := c.yagkclient.LoginUser(ctx, &pb.User{Email: user.Email, Password: user.Password})
	if err != nil {
		return "", err
	}
	c.token = resp.Token
	return resp.Token, nil
}

func (c *Client) CheckSync(email string) (time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := c.yagkclient.CheckSync(ctx, &pb.CheckSyncRequest{Email: email})
	if err != nil {
		return time.Time{}, err
	}
	lastSync := time.Unix(resp.LastSync.Seconds, 0)
	resp.GetLastSync()
	if err != nil {
		return time.Time{}, err
	}
	return lastSync, nil
}

func (c *Client) GetData(email string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := c.yagkclient.GetData(ctx, &emptypb.Empty{})
	if err != nil || len(resp.Data) == 0 {
		return nil, err
	}
	return resp.Data, nil
}

func (c *Client) SendData(email string, data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	secrets := &pb.Secrets{Data: data}
	resp, err := c.yagkclient.SetData(ctx, secrets)
	resp.GetLastSync()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) AuthInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		ctx = metadata.AppendToOutgoingContext(ctx, "Bearer", c.token)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
