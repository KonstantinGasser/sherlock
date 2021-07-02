package internal

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/KonstantinGasser/required"
)

const (
	prettyDateLayout = "Monday, 02. January 2006"
)

var (
	ErrAccountExists    = fmt.Errorf("account for group already exists")
	ErrNoSuchAccount    = fmt.Errorf("account not found")
	ErrInvalidGroupName = fmt.Errorf("group name must be a consecutive string")
)

// Group groups Accounts
type Group struct {
	GID      string     `json:"name" required:"yes"`
	Accounts []*Account `json:"accounts"`
}

func NewGroup(name string) (*Group, error) {
	g := Group{
		GID:      name,
		Accounts: make([]*Account, 0),
	}
	if err := g.valid(); err != nil {
		return nil, err
	}
	return &g, nil
}

// append appends an account to a group if it does not already exists
func (g *Group) append(account *Account) error {
	if ok := g.exists(account); ok {
		return ErrAccountExists
	}
	g.Accounts = append(g.Accounts, account)
	return nil
}

func (g Group) lookup(accountName string) (*Account, error) {
	for _, a := range g.Accounts {
		if a.Name == accountName {
			return a, nil
		}
	}
	return nil, ErrNoSuchAccount
}

// delete deletes a given account from the group, returns an ErrNoSuchAccount
// if account not present
func (g *Group) delete(account string) error {
	var offset *int
	for i, a := range g.Accounts {
		if a.Name == account {
			offset = &i
		}
	}
	if offset == nil {
		return ErrNoSuchAccount
	}

	g.Accounts = append(g.Accounts[:*offset], g.Accounts[*offset+1:]...)
	return nil
}

// exists checks an account is already present in the group
// the Account.Name and Account.Tag field build the pk for an Account
func (g Group) exists(account *Account) bool {
	for _, a := range g.Accounts {
		if account.Name == a.Name && account.Tag == a.Tag {
			return true
		}
	}
	return false
}

func (g Group) serizalize() ([]byte, error) {
	return json.Marshal(g)
}

func (g Group) valid() error {
	if err := required.Atomic(&g); err != nil {
		return ErrMissingValues
	}
	if set := strings.Split(g.GID, " "); len(set) != 1 {
		return ErrInvalidGroupName
	}
	return nil
}

// Table builds the Group in such a way that it can be consumed by the tablewriter.Table
func (g Group) Table(filter ...func(*Account) bool) [][]string {
	var accounts [][]string

skipp:
	for _, item := range g.Accounts {
		for _, f := range filter {
			if !f(item) {
				continue skipp
			}
		}
		accounts = append(accounts, []string{
			g.GID,
			item.Name,
			strings.Join([]string{"#", item.Tag}, ""),
			item.CreatedOn.Format(prettyDateLayout),
			item.UpdatedOn.Format(prettyDateLayout),
		})
	}
	return accounts
}

func FilterByTag(tag string) func(*Account) bool {
	return func(a *Account) bool {
		if len(tag) == 0 {
			return true
		}
		return a.Tag == tag
	}
}
