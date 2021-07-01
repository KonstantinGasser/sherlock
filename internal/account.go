package internal

import (
	"fmt"

	"github.com/KonstantinGasser/sherlock/internal/security"
)

var (
	ErrInsecurePassword = fmt.Errorf("provided password is insecure (use --insecure to ignore this message)")
)

type Account struct {
	// reference to a Group if set
	Group    string `proto:"gid"`
	Name     string `proto:"name"`
	Password string `proto:"password"`
	Desc     string `proto:"desc"`
}

// NewAccount creates a new Account and if insecure=false checks the password strength
// returning an err if strength security.Low
func NewAccount(gid, name, password, desc string, insecure bool) (*Account, error) {
	a := Account{
		Group:    gid,
		Name:     name,
		Password: password,
		Desc:     desc,
	}
	if insecure {
		return &a, nil
	}
	fmt.Println(a.secure())
	if level := a.secure(); level == security.Low {
		return nil, ErrInsecurePassword
	}
	return &a, nil
}

// secure checks the Accounts on how secure it is
func (a Account) secure() int {
	return security.PasswordStrength(a.Password)
}
