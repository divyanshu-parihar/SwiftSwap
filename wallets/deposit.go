package wallets

type DespositWallet struct {
	Crypto string
	Addr   string
	Chain  string
}

func NewDepositWallets() []*DespositWallet {
	return []*DespositWallet{
		&DespositWallet{
			Crypto: "USDT",
			Addr:   "TDdUG5jw9Afje8FfkWaxFetiyJX6reWgtH",
			Chain:  "TRC20",
		},
		&DespositWallet{
			Crypto: "BTC",
			Addr:   "3FYnYbSTBVVVXh1f23E6Ufijc57zHpCjMG",
			Chain:  "BTC",
		},
		&DespositWallet{
			Crypto: "ETH",
			Addr:   "0x357d5e6a96db33825b448603a136f92c8e244b0f",
			Chain:  "ERC20",
		},
		&DespositWallet{
			Crypto: "ETH",
			Addr:   "0x357d5e6a96db33825b448603a136f92c8e244b0f",
			Chain:  "BASE",
		},
	}
}
