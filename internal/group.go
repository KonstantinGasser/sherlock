package internal

import (
	"encoding/json"
	"fmt"
)

var (
	ErrAccountExists = fmt.Errorf("account for group already exists")
)

// Group groups Accounts
type Group struct {
	GID      string     `json:"name"`
	Accounts []*Account `json:"accounts"`
}

func NewGroup(name string) *Group {
	return &Group{
		GID:      name,
		Accounts: make([]*Account, 0),
	}
}

func (g *Group) append(accounts ...*Account) error {
	for _, account := range accounts {
		if ok := g.exists(account); ok {
			return ErrAccountExists
		}
	}
	g.Accounts = append(g.Accounts, accounts...)
	return nil
}

func (g Group) exists(account *Account) bool {
	for _, a := range g.Accounts {
		if account.Name == a.Name {
			return true
		}
	}
	return false
}

func (g Group) serizalize() ([]byte, error) {
	return json.Marshal(g)
}
