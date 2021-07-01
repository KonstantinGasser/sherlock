package internal

// Group groups Accounts
type Group struct {
	GID      string
	Desc     string
	Accounts []*Account
}

func NewGroup(name string, desc string) *Group {
	return &Group{
		GID:      name,
		Desc:     desc,
		Accounts: make([]*Account, 0),
	}
}
