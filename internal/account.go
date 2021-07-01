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
	Tag       string    `json:"tag"`
	CreatedOn time.Time `json:"created_on"`
}

// NewAccount creates a new Account and if insecure=false checks the password strength
// returning an err if strength security.Low
func NewAccount(name, password, tag string, insecure bool) (*Account, error) {
	a := Account{
		Name:      name,
		Password:  password,
		CreatedOn: time.Now(),
		Tag:       tag,
	}
	if insecure {
		return &a, nil
	}
	if level := a.secure(); level == security.Low {
		return nil, ErrInsecurePassword
	}
	return &a, nil
}

// secure checks the Accounts on how secure it is
func (a Account) secure() int {
	return security.PasswordStrength(a.Password)
}
