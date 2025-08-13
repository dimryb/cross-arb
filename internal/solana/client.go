package solana

import (
	"context"
	"fmt"
	"strings"

	i "github.com/dimryb/cross-arb/internal/interface"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	rpctx "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

// Client отвечает за взаимодействие с блокчейном Solana.
type Client struct {
	rpcClient *rpc.Client
	wsClient  *ws.Client
	logger    i.Logger
}

// NewSolanaClient создает клиент для работы с Solana RPC.
func NewSolanaClient(logger i.Logger, rpcURL string) (*Client, error) {
	wsURL := strings.Replace(rpcURL, "https://", "wss://", 1)
	wsURL = strings.Replace(wsURL, "http://", "ws://", 1)

	wsClient, err := ws.Connect(context.Background(), wsURL)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к Solana WebSocket: %w", err)
	}

	return &Client{
		rpcClient: rpc.New(rpcURL),
		wsClient:  wsClient,
		logger:    logger,
	}, nil
}

// SendAndConfirmTransaction отправляет и подтверждает транзакцию.
func (c *Client) SendAndConfirmTransaction(ctx context.Context, tx *solana.Transaction) (solana.Signature, error) {
	sig, err := rpctx.SendAndConfirmTransaction(ctx, c.rpcClient, c.wsClient, tx)
	if err != nil {
		return solana.Signature{}, err
	}

	c.logger.Debug("Транзакция успешно отправлена и подтверждена", "сигнатура", sig.String())
	return sig, nil
}

// GetBalance получает баланс аккаунта.
func (c *Client) GetBalance(ctx context.Context, account solana.PublicKey) (uint64, error) {
	balance, err := c.rpcClient.GetBalance(ctx, account, rpc.CommitmentFinalized)
	if err != nil {
		return 0, err
	}
	return balance.Value, nil
}

// Close закрывает соединения.
func (c *Client) Close() {
	c.wsClient.Close()
}
