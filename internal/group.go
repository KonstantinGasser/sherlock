package internal

import (
	"encoding/json"
	"fmt"
)

const (
	prettyDateLayout = "Monday, 02. January 2006"
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

func (g Group) Table() [][]string {
	var accounts = make([][]string, len(g.Accounts))
	for i, item := range g.Accounts {
		accounts[i] = []string{
			g.GID,
			item.Name,
			item.Desc,
			item.CreatedOn.Format(prettyDateLayout),
		}
	}
	return accounts
}
