package terminal

import (
	"fmt"
	"syscall"

	"github.com/enescakir/emoji"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
)

func Success(format string, a ...interface{}) {
	pretty(color.FgGreen, emoji.Emoji(emoji.RaisingHands.String()), format, a...)
}

func Error(format string, a ...interface{}) {
	pretty(color.FgRed, emoji.ExclamationMark, format, a...)
}

func ReadPassword(format string) (string, error) {
	prettyNoNewLine(color.FgGreen, emoji.Locked, format)
	b, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Print("\n")
	return string(b), nil
}

// pretty combines the colors and emojis and outputs a formatted string to the
// cli
func pretty(c color.Attribute, e emoji.Emoji, f string, a ...interface{}) {
	_, _ = color.New(c).Printf(fmt.Sprintf("%v %s\n", e, f), a...)
}

// prettyNoNewLine combines the colors and emojis and outputs a formatted string to the
// cli. does not add a \n to the format string
func prettyNoNewLine(c color.Attribute, e emoji.Emoji, f string, a ...interface{}) {
	_, _ = color.New(c).Printf(fmt.Sprintf("%v %s", e, f), a...)
}
