package investgo

import (
	"context"

	pb "github.com/strukchinsky/tbank-invest-go-sdk"
)

func (c *Client) GetLastPrices(ctx context.Context, id []string, lastPriceType *pb.LastPriceType) (map[string]*pb.Quotation, error) {
	client := pb.NewMarketDataServiceClient(c.conn)

	prices, err := client.GetLastPrices(ctx, &pb.GetLastPricesRequest{
		InstrumentId:  id,
		LastPriceType: *lastPriceType,
	})
	if err != nil {
		return nil, err
	}

	priceById := make(map[string]*pb.Quotation)
	for _, price := range prices.LastPrices {
		priceById[price.InstrumentUid] = price.Price
	}

	return priceById, nil
}

func (c *Client) GetLastPrice(ctx context.Context, id string, lastPriceType *pb.LastPriceType) (*pb.Quotation, error) {
	client := pb.NewMarketDataServiceClient(c.conn)

	prices, err := client.GetLastPrices(ctx, &pb.GetLastPricesRequest{
		InstrumentId:  []string{id},
		LastPriceType: *lastPriceType,
	})
	if err != nil {
		return nil, err
	}

	return prices.LastPrices[0].Price, nil
}
