package investgo

import (
	"context"
	"crypto/tls"
	"fmt"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	token string
	conn  *grpc.ClientConn
}

func EnrichErrorWithMessageAndTrackingIdUnaryClientInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	var trailer metadata.MD
	err := invoker(ctx, method, req, reply, cc, append(opts, grpc.Trailer(&trailer))...)

	if err != nil {
		err = fmt.Errorf("%s (%w; x-tracking-id = %s)", MessageFromMetadata(trailer), err, TrackingIdFromMetadata(trailer))
	}

	return err
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
		grpc.WithChainUnaryInterceptor(EnrichErrorWithMessageAndTrackingIdUnaryClientInterceptor),
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
