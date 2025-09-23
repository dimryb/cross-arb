package jupiter

// MintPair хранит адреса SPL-токенов торговой пары.
type MintPair struct {
	// InputMint — адрес BASE-токена в CEX-биржах и INPUT в DEX-биржах.
	// Например: для пары SOL/USDC это будет SOL
	InputMint string
	// OutputMint — адрес QUOTE-токена в CEX-биржах и OUTPUT в DEX-биржах.
	// Например: для пары SOL/USDC это будет USDC
	OutputMint string
}
