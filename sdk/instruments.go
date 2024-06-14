package investgo

import (
	"context"

	pb "github.com/floatdrop/tbank-invest-go-sdk"
)

func (c *Client) Shares(ctx context.Context, status pb.InstrumentStatus, exchange pb.InstrumentExchangeType) ([]*pb.Share, error) {
	client := pb.NewInstrumentsServiceClient(c.conn)

	response, err := client.Shares(ctx, &pb.InstrumentsRequest{
		InstrumentStatus:   &status,
		InstrumentExchange: &exchange,
	})
	if err != nil {
		return nil, err
	}

	return response.Instruments, nil
}

func ByFigi(figi string) pb.InstrumentRequest {
	return pb.InstrumentRequest{
		IdType: pb.InstrumentIdType_INSTRUMENT_ID_TYPE_FIGI,
		Id:     figi,
	}
}

func ByTicker(ticker string, classCode string) pb.InstrumentRequest {
	return pb.InstrumentRequest{
		IdType:    pb.InstrumentIdType_INSTRUMENT_ID_TYPE_TICKER,
		Id:        ticker,
		ClassCode: &classCode,
	}
}

func ByUid(uid string) pb.InstrumentRequest {
	return pb.InstrumentRequest{
		IdType: pb.InstrumentIdType_INSTRUMENT_ID_TYPE_UID,
		Id:     uid,
	}
}

func ByPositionUid(positionUid string) pb.InstrumentRequest {
	return pb.InstrumentRequest{
		IdType: pb.InstrumentIdType_INSTRUMENT_ID_TYPE_POSITION_UID,
		Id:     positionUid,
	}
}

func (c *Client) ShareBy(ctx context.Context, request *pb.InstrumentRequest) (*pb.Share, error) {
	client := pb.NewInstrumentsServiceClient(c.conn)

	response, err := client.ShareBy(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Instrument, nil
}
