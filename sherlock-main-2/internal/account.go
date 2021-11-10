package internal

import (
	"fmt"
	"strings"
	"time"

	"github.com/KonstantinGasser/required"
	"github.com/KonstantinGasser/sherlock/security"
	"github.com/m1/go-generate-password/generator"
)

var (
	ErrInsecurePassword         = fmt.Errorf("provided password is insecure (use --insecure to ignore this message)")
	ErrInvalidAccountName       = fmt.Errorf("account name must be a consecutive string")
	ErrMissingValues            = fmt.Errorf("account is missing required values")
	ErrInvalidAccountNameSymbol = fmt.Errorf("account name invalid. Please avoid using '@' character")
)

// fieldUpdate is a function which can alter the fields of
// an account
type fieldUpdate func(*account) error

type account struct {
	Name      string    `json:"name" required:"yes"`
	Password  string    `json:"password" required:"yes"`
	Tag       string    `json:"tag"`
	CreatedOn time.Time `json:"created_on" required:"yes"`
	UpdatedOn time.Time `json:"updated_on"`
}

// NewAccount creates a new Account and if insecure=false checks the password strength
// returning an err if strength security.Low
func NewAccount(query, password, tag string, insecure bool) (*account, error) {
	_, acc, err := SplitQuery(query)
	if err != nil {
		return nil, err
	}
	a := account{
		Name:      acc,
		Password:  password,
		CreatedOn: time.Now(),
		UpdatedOn: time.Now(),
		Tag:       tag,
	}
	if err := a.valid(); err != nil {
		return nil, err
	}

	if insecure {
		return &a, nil
	}
	if err := a.secure(); err != nil {
		return nil, err
	}
	return &a, nil
}

func (a account) valid() error {
	if err := required.Atomic(&a); err != nil {
		return ErrMissingValues
	}
	if set := strings.Split(a.Name, " "); len(set) > 1 {
		return ErrInvalidAccountName
	}
	return nil
}

func updateFieldName(name string) fieldUpdate {
	return func(a *account) error {
		a.Name = strings.TrimSpace(name)
		return nil
	}
}

func updateFieldPassword(password string, insecure bool) fieldUpdate {
	return func(a *account) error {
		a.Password = strings.TrimSpace(password)
		if insecure {
			a.UpdatedOn = time.Now()
			return nil
		}
		if err := a.secure(); err != nil {
			return err
		}
		return nil
	}
}

func updateFieldTag(tag string) fieldUpdate {
	return func(a *account) error {
		a.Tag = strings.TrimSpace(tag)
		return nil
	}
}

func (a *account) update(opt fieldUpdate) error {
	if err := opt(a); err != nil {
		return err
	}
	a.UpdatedOn = time.Now()
	return nil
}

// secure checks the Accounts on how secure it is
func (a account) secure() error {
	return security.PasswordStrength(a.Password)
}

func AutoGeneratePassword(passwordLength int) (string, error) {
	config := generator.Config{
		Length:                     passwordLength,
		IncludeSymbols:             true,
		IncludeNumbers:             true,
		IncludeLowercaseLetters:    true,
		IncludeUppercaseLetters:    true,
		ExcludeSimilarCharacters:   true,
		ExcludeAmbiguousCharacters: true,
	}
	g, _ := generator.New(&config)

	pwd, err := g.Generate()
	return *pwd, err
}
