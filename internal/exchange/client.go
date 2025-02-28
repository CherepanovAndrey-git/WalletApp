package exchange

import (
	"context"
	"fmt"
	"time"

	pb "github.com/CherepanovAndrey-git/gw-exchange-grpc/proto/exchange"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.ExchangeServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	return &Client{
		conn:   conn,
		client: pb.NewExchangeServiceClient(conn),
	}, nil
}

func (c *Client) GetRates(ctx context.Context) (map[string]float32, error) {

	resp, err := c.client.GetExchangeRates(ctx, &pb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("failed to get rates: %w", err)
	}
	return resp.Rates, nil
}

func (c *Client) GetRate(from, to string) (float32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := c.client.GetExchangeRateForCurrency(ctx, &pb.CurrencyRequest{
		FromCurrency: from,
		ToCurrency:   to,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to get rate: %w", err)
	}
	return resp.Rate, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
