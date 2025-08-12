package interfaces

import "github.com/gagliardetto/solana-go"

type TransactionSigner interface {
	PublicKey() solana.PublicKey
	SignTransaction(tx *solana.Transaction) error
}
