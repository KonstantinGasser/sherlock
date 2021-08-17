package internal

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/KonstantinGasser/required"
	"github.com/KonstantinGasser/sherlock/security"
)

const (
	defaultGroupName = "default"
	prettyDateLayout = "Monday, 02. January 2006"
)

var (
	ErrAccountExists          = fmt.Errorf("account for group already exists")
	ErrNoSuchAccount          = fmt.Errorf("account not found")
	ErrInvalidGroupName       = fmt.Errorf("group name must be a consecutive string")
	ErrInvalidGroupNameSymbol = fmt.Errorf("group name invalid. Please avoid using '@' character")
)

// Group groups Accounts
type group struct {
	GID      string     `json:"name" required:"yes"`
	Accounts []*account `json:"accounts"`
}

func newDefaultGroup() *group {
	return &group{
		GID:      defaultGroupName,
		Accounts: make([]*account, 0),
	}
}

func NewGroup(name string) (*group, error) {
	g := group{
		GID:      name,
		Accounts: make([]*account, 0),
	}
	if err := g.valid(); err != nil {
		return nil, err
	}
	return &g, nil
}

// append appends an account to a group if it does not already exists
func (g *group) append(account *account) error {
	if ok := g.exists(account.Name); ok {
		return ErrAccountExists
	}
	g.Accounts = append(g.Accounts, account)
	return nil
}

func (g group) lookup(accountName string) (*account, error) {
	for _, a := range g.Accounts {
		if a.Name == accountName {
			return a, nil
		}
	}
	return nil, ErrNoSuchAccount
}

// delete deletes a given account from the group, returns an ErrNoSuchAccount
// if account not present
func (g *group) delete(account string) error {
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
// using the account.Name as a pk
func (g group) exists(name string) bool {
	for _, a := range g.Accounts {
		if name == a.Name {
			return true
		}
	}
	return false
}

func (g group) serizalize() ([]byte, error) {
	return json.Marshal(g)
}

func (g group) valid() error {
	if err := required.Atomic(&g); err != nil {
		return ErrMissingValues
	}
	if set := strings.Split(g.GID, " "); len(set) != 1 {
		return ErrInvalidGroupName
	}
	return nil
}

// secure evaluates the password strength of the group password
func (g group) secure(groupKey string) error {
	return security.PasswordStrength(groupKey)
}

// Table builds the Group in such a way that it can be consumed by the tablewriter.Table
func (g group) Table(filter ...func(*account) bool) [][]string {
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

func FilterByTag(tag string) func(*account) bool {
	return func(a *account) bool {
		if len(tag) == 0 {
			return true
		}
		return a.Tag == tag
	}
}
