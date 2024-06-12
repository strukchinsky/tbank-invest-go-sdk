package investgo

import (
	"context"
	"slices"
	"time"

	pb "github.com/floatdrop/tbank-invest-go-sdk"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetCandles splits interval into day chunks and loads HistoricCandle for it.
func (c *Client) GetCandles(ctx context.Context, figi string, from time.Time, to time.Time, interval pb.CandleInterval) ([]*pb.HistoricCandle, error) {
	client := pb.NewMarketDataServiceClient(c.conn)

	nextDay := from.Add(time.Hour * 24)
	result := make([]*pb.HistoricCandle, 0)

	for {
		candles, err := client.GetCandles(ctx, &pb.GetCandlesRequest{
			Figi:     &figi,
			From:     timestamppb.New(from),
			To:       timestamppb.New(nextDay),
			Interval: interval,
		})
		if err != nil {
			return nil, err
		}

		result = append(result, candles.Candles...)

		if nextDay == to {
			break
		}

		from = nextDay
		nextDay = slices.MinFunc([]time.Time{from.Add(time.Hour * 24), to}, time.Time.Compare)
	}
	return result, nil
}
