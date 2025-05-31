package investgo

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/strukchinsky/tbank-invest-go-sdk"
)

func (c *Client) PostOrder(ctx context.Context, request *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	client := pb.NewOrdersServiceClient(c.conn)
	return client.PostOrder(ctx, request)
}

func (c *Client) BuyMarket(ctx context.Context, accountId string, instrumentId string, quantity int64) (*pb.PostOrderResponse, error) {
	client := pb.NewOrdersServiceClient(c.conn)

	request := &pb.PostOrderRequest{
		InstrumentId: instrumentId,
		Quantity:     quantity,
		Direction:    pb.OrderDirection_ORDER_DIRECTION_BUY,
		AccountId:    accountId,
		OrderType:    pb.OrderType_ORDER_TYPE_MARKET,
		OrderId:      uuid.NewString(),
	}

	return client.PostOrder(ctx, request)
}

func (c *Client) SellMarket(ctx context.Context, accountId string, instrumentId string, quantity int64) (*pb.PostOrderResponse, error) {
	client := pb.NewOrdersServiceClient(c.conn)

	request := &pb.PostOrderRequest{
		InstrumentId: instrumentId,
		Quantity:     quantity,
		Direction:    pb.OrderDirection_ORDER_DIRECTION_SELL,
		AccountId:    accountId,
		OrderType:    pb.OrderType_ORDER_TYPE_MARKET,
		OrderId:      uuid.NewString(),
	}

	return client.PostOrder(ctx, request)
}

func (c *Client) GetMaxLots(ctx context.Context, request *pb.GetMaxLotsRequest) (*pb.GetMaxLotsResponse, error) {
	client := pb.NewOrdersServiceClient(c.conn)
	return client.GetMaxLots(ctx, request)
}
