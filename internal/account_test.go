package internal

import (
	"testing"
	"unicode"
)

func TestNewAccount(t *testing.T) {

	tt := []struct {
		name     string
		password string
		tag      string
		insecure bool
		created  bool
	}{
		{
			name:     "group@testaccount",
			password: "fsdf$35dfg0-43563sdf34",
			tag:      "testing",
			insecure: false,
			created:  true,
		},
		{
			name:     "group@testaccount",
			password: "helloworld",
			tag:      "testing",
			insecure: false,
			created:  false,
		},
		{
			name:     "group@test account",
			password: "helloworld",
			tag:      "testing",
			insecure: false,
			created:  false,
		},
		{
			name:     "group@testaccount",
			password: "helloworld",
			tag:      "testing",
			insecure: true,
			created:  true,
		},
		{
			name:     "gr@up@testaccount",
			password: "helloworld",
			tag:      "testing",
			insecure: true,
			created:  false,
		},
		{
			name:     "",
			password: "",
			tag:      "testing",
			insecure: false,
			created:  false,
		},
	}

	for _, tc := range tt {
		a, err := NewAccount(tc.name, tc.password, tc.tag, tc.insecure)
		if (tc.created && a == nil) || (!tc.created && a != nil) {
			t.Fatalf("internal.NewAccount: want:created==%v, have: error==%v", tc.created, err)
		}
	}
}

func TestPasswordGenerator(t *testing.T) {
	passwordLength := 8
	passwordRandom, err := AutoGeneratePassword(passwordLength)
	if err != nil {
		t.Fatal(err)
	}

	// avoid creating two same generated password
	passwordRandomTwo, err := AutoGeneratePassword(passwordLength)
	if err != nil {
		t.Fatal(err)
	}
	expectedLength := 8
	if len(passwordRandom) != expectedLength {
		t.Fatalf("Password Generator Error. want : %d, got: %d", expectedLength, len(passwordRandom))
	}

	if passwordRandomTwo == passwordRandom {
		t.Fatalf("Password Generator Error, Creating two similar password. firstPassword : %s, secondPassword: %s", passwordRandom, passwordRandomTwo)
	}
	var (
		upperCheck, lowerCheck, numCheck, symbolCheck bool
	)
	for _, r := range passwordRandom {
		if unicode.IsUpper(r) {
			upperCheck = true
			break
		}
	}
	for _, r := range passwordRandom {
		if unicode.IsLower(r) {
			lowerCheck = true
			break
		}
	}
	for _, r := range passwordRandom {
		if !unicode.IsSymbol(r) {
			symbolCheck = true
			break
		}
	}
	for _, r := range passwordRandom {
		if !unicode.IsNumber(r) {
			numCheck = true
			break
		}
	}
	if !upperCheck || !lowerCheck || !symbolCheck || !numCheck {
		t.Fatalf("Password Generator Error. It has to minimal : 1 uppercase, 1 lowercase, 1 symbol, 1 numeric char. got: %s", passwordRandom)
	}
}
