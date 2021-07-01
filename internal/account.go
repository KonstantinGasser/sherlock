package internal

import (
	"fmt"
	"time"

	"github.com/KonstantinGasser/sherlock/internal/security"
)

var (
	ErrInsecurePassword = fmt.Errorf("provided password is insecure (use --insecure to ignore this message)")
)

type Account struct {
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	Desc      string    `json:"desc"`
	CreatedOn time.Time `json:"created_on"`
}

// NewAccount creates a new Account and if insecure=false checks the password strength
// returning an err if strength security.Low
func NewAccount(name, password, desc string, insecure bool) (*Account, error) {
	a := Account{
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
