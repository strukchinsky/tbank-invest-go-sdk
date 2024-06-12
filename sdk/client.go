package investgo

import (
	"crypto/tls"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

type Client struct {
	token string
	conn  *grpc.ClientConn
}

// NewClient creates SDK client that will be used for all communications with investapi
//
// By default it configures default TLS transport credentials and oauth2.StaticTokenSource for PerRPCCredentials
func NewClient(target string, token string, opts ...grpc.DialOption) (*Client, error) {
	defaultOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
		grpc.WithPerRPCCredentials(oauth.TokenSource{
			TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
		}),
	}

	conn, err := grpc.NewClient(
		target,
		append(defaultOpts, opts...)...,
	)
	if err != nil {
		return nil, err
	}

	return &Client{token, conn}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
