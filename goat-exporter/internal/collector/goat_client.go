package collector

import (
	"context"
	"fmt"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
)

const dialTimeout = 10 * time.Second

// GoatClient wraps go-ethereum's ethclient with helpers tailored for the exporter.
type GoatClient struct {
	client *ethclient.Client
}

// NewGoatClient connects to the Goat network RPC endpoint using go-ethereum's ethclient.
func NewGoatClient(endpoint string) (*GoatClient, error) {
	if strings.TrimSpace(endpoint) == "" {
		return nil, fmt.Errorf("endpoint cannot be empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()

	client, err := ethclient.DialContext(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("connect to rpc endpoint: %w", err)
	}

	return &GoatClient{client: client}, nil
}

func (c *GoatClient) BlockNumber(ctx context.Context) (uint64, error) {
	return c.client.BlockNumber(ctx)
}

func (c *GoatClient) ChainID(ctx context.Context) (uint64, error) {
	chainID, err := c.client.ChainID(ctx)
	if err != nil {
		return 0, err
	}
	return chainID.Uint64(), nil
}

func (c *GoatClient) SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error) {
	return c.client.SyncProgress(ctx)
}

func (c *GoatClient) Close() {
	c.client.Close()
}
