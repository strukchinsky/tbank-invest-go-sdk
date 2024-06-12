package investgo

import (
	"context"

	pb "github.com/floatdrop/tbank-invest-go-sdk"
)

func (c *Client) Shares(ctx context.Context, status pb.InstrumentStatus) ([]*pb.Share, error) {
	client := pb.NewInstrumentsServiceClient(c.conn)

	response, err := client.Shares(ctx, &pb.InstrumentsRequest{InstrumentStatus: &status})
	if err != nil {
		return nil, err
	}

	return response.Instruments, nil
}
