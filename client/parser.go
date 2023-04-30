package client

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/c-bata/go-prompt"
)

var Parser = &TermParser{}

type TermParser struct {
	prefix_string string
	cur_buffer    string
	p             *prompt.Prompt
}

func executor(in string) {
	fmt.Printf("Your input: %s, state: %d\n", in, GameState.phase)
}

func completer(in prompt.Document) []prompt.Suggest {
	Parser.cur_buffer = in.Text

	s := []prompt.Suggest{
		{Text: "users", Description: "Store the username and age"},
		{Text: "articles", Description: "Store the article text posted by user"},
		{Text: "comments", Description: "Store the text commented to articles"},
		{Text: "groups", Description: "Combine users with specific rules"},
		{Text: "exit", Description: "Combine users with specific rules"},
	}

	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func exitChecker(in string, breakline bool) bool {
	return in == "exit" && breakline
}

func handleExit() {
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
}

func (tp *TermParser) GetCurBuf() string {
	return tp.cur_buffer
}

func (tp *TermParser) GetPrefixString() string {
	green := "\033[32m"
	reset := "\033[0m"
	return fmt.Sprintf("%s%s%s", green, tp.prefix_string, reset)
}

func (tp *TermParser) Run() {
	defer handleExit()
	tp.p.Run()
}

func (tp *TermParser) Init() {
	tp.cur_buffer = ""
	tp.prefix_string = ">>> "
	tp.p = prompt.New(
		executor,
		completer,
		prompt.OptionSetExitCheckerOnInput(exitChecker),
		prompt.OptionPrefix(tp.prefix_string),
		prompt.OptionPrefixTextColor(prompt.DarkGreen),
		prompt.OptionPreviewSuggestionTextColor(prompt.White),
		prompt.OptionSelectedSuggestionTextColor(prompt.DarkGreen),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionDescriptionBGColor(prompt.White),
	)
}
