package investgo

import (
	"context"

	pb "github.com/strukchinsky/tbank-invest-go-sdk"
)

func (c *Client) PostOrder(ctx context.Context, request *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	client := pb.NewOrdersServiceClient(c.conn)
	return client.PostOrder(ctx, request)
}

func (c *Client) GetMaxLots(ctx context.Context, request *pb.GetMaxLotsRequest) (*pb.GetMaxLotsResponse, error) {
	client := pb.NewOrdersServiceClient(c.conn)
	return client.GetMaxLots(ctx, request)
}
