package grpcclient

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/Spear5030/yagophkeeper/internal/domain"
	"github.com/Spear5030/yagophkeeper/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"os"
	"time"
)

type Client struct {
	conn       *grpc.ClientConn
	yagkclient pb.YaGophKeeperClient
	token      string
}

func New(addr string, cert string, token string) *Client {
	var c Client
	tlsCredentials, err := loadTLSCredentials(cert)
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(tlsCredentials),
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
	return lastSync, nil
}

func (c *Client) GetData() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := c.yagkclient.GetData(ctx, &emptypb.Empty{})
	if err != nil || len(resp.Data) == 0 {
		return nil, err
	}
	return resp.Data, nil
}

func (c *Client) SendData(data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	secrets := &pb.Secrets{Data: data}
	_, err := c.yagkclient.SetData(ctx, secrets)
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

func loadTLSCredentials(cert string) (credentials.TransportCredentials, error) {
	pemServerCA, err := os.ReadFile(cert)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	config := &tls.Config{
		RootCAs:            certPool,
		InsecureSkipVerify: true, //
	}

	return credentials.NewTLS(config), nil
}
