package bot

type AddWalletsForm struct {
	User    int64
	Network string
	Coin    string
	Address string
	Memo    string
	Step    int // 0: User, 1: Network, 2: Address, 3: Memo
}
type QuickTransaction struct {
	User    string
	Network string
	Coin    string
	Address string
	Memo    string

	Step int // 0: User, 1: Network, 2: Address, 3: Memo
}

type UserTrustedWalletForm struct {
	User    int64
	Address string
}
