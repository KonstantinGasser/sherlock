package terminal

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/enescakir/emoji"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/crypto/ssh/terminal"
)

func Success(format string, a ...interface{}) {
	pretty(color.FgGreen, emoji.Emoji(emoji.RaisingHands.String()), format, a...)
}

func Info(format string, a ...interface{}) {
	pretty(color.FgHiBlue, emoji.Emoji(emoji.BackhandIndexPointingRight.String()), format, a...)
}

func Error(format string, a ...interface{}) {
	pretty(color.FgRed, emoji.ExclamationMark, format, a...)
}

func ReadPassword(format string) (string, error) {
	prettyNoNewLine(color.FgHiBlue, emoji.Locked, format)
	b, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Print("\n")
	return string(b), nil
}

// YesNo prompts the user with a confirm dialog. in every case except for "y"
// (lowercase y) the return will be false
func YesNo(format string) bool {
	r := bufio.NewReader(os.Stdin)
	prettyNoNewLine(color.FgRed, emoji.FaceWithMonocle, format)
	input, _ := r.ReadString('\n')

	switch strings.TrimSuffix(input, "\n") {
	case "y":
		return true
	default:
		return false
	}
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

var bgC = []int{
	tablewriter.BgBlueColor,
	tablewriter.BgMagentaColor,
	tablewriter.BgGreenColor,
	tablewriter.BgHiYellowColor,
}

func ToTable(header []string, rows [][]string, opts ...func(*tablewriter.Table)) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(padding(header))
	buildHeader(table, header)

	for _, opt := range opts {
		opt(table)
	}
	table.AppendBulk(rows)
	table.Render()
}

func buildHeader(t *tablewriter.Table, h []string) {
	colors := make([]tablewriter.Colors, len(h))
	for i := 0; i < len(h); i++ {
		colors[i] = tablewriter.Colors{tablewriter.Bold, bgC[i%len(h)]}
	}
	t.SetHeaderColor(colors...)
}

func padding(h []string) []string {
	for i, v := range h {
		h[i] = " " + v + " "
	}
	return h
}

// TableWithCellMerge apply tablewriter.SetAuthMergeCellsByColumnIndex to the
// table instance and enables tablewriter.SetRowLine.
// Allows to group rows by a column index
func TableWithCellMerge(mergeByIndex int) func(*tablewriter.Table) {
	return func(t *tablewriter.Table) {
		var index = mergeByIndex
		if mergeByIndex > t.NumLines() {
			index = 0
		}
		t.SetAutoMergeCellsByColumnIndex([]int{index})
		t.SetRowLine(true)
	}
}
