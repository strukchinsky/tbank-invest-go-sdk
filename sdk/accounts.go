package investgo

import (
	"context"

	pb "github.com/strukchinsky/tbank-invest-go-sdk"
)

func (c *Client) GetAccounts(ctx context.Context) ([]*pb.Account, error) {
	client := pb.NewUsersServiceClient(c.conn)

	in := pb.GetAccountsRequest{
		Status: pb.AccountStatus_ACCOUNT_STATUS_ALL.Enum(),
	}

	response, err := client.GetAccounts(ctx, &in)
	if err != nil {
		return nil, err
	}

	return response.Accounts, nil
}
