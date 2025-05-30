package investgo

import (
	"context"

	pb "github.com/strukchinsky/tbank-invest-go-sdk"
)

func (c *Client) GetPortfolio(ctx context.Context, accountId string, currency pb.PortfolioRequest_CurrencyRequest) (*pb.PortfolioResponse, error) {
	client := pb.NewOperationsServiceClient(c.conn)

	in := pb.PortfolioRequest{
		AccountId: accountId,
		Currency:  &currency,
	}

	return client.GetPortfolio(ctx, &in)
}
