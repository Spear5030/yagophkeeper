package grpcclient

import (
	"context"
	"github.com/Spear5030/yagophkeeper/internal/domain"
	"github.com/Spear5030/yagophkeeper/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

type Client struct {
	conn       *grpc.ClientConn
	yagkclient pb.YaGophKeeperClient
	token      string
}

func New(addr string) *Client {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	return &Client{
		conn:       conn,
		yagkclient: pb.NewYaGophKeeperClient(conn),
		token:      "",
	}

}

func (c *Client) RegisterUser(user domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := c.yagkclient.RegisterUser(ctx, &pb.User{Email: user.Email, Password: user.Password})
	if err != nil {
		return err
	}
	c.token = resp.Token
	return nil
}

func (c *Client) LoginUser(user domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := c.yagkclient.LoginUser(ctx, &pb.User{Email: user.Email, Password: user.Password})
	if err != nil {
		return err
	}
	c.token = resp.Token
	return nil
}
