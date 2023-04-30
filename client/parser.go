package client

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/c-bata/go-prompt"
)

var Parser = makeParser()

type TermParser struct {
	prefix_string string
	p             *prompt.Prompt
}

func (tp *TermParser) executor(in string) {
	fmt.Println("Your input: " + in)
}

var cur_buffer string = ""

func (tp *TermParser) completer(in prompt.Document) []prompt.Suggest {
	cur_buffer = in.Text

	s := []prompt.Suggest{
		{Text: "users", Description: "Store the username and age"},
		{Text: "articles", Description: "Store the article text posted by user"},
		{Text: "comments", Description: "Store the text commented to articles"},
		{Text: "groups", Description: "Combine users with specific rules"},
		{Text: "exit", Description: "Combine users with specific rules"},
	}

	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func (tp *TermParser) exitChecker(in string, breakline bool) bool {
	return in == "exit" && breakline
}

func handleExit() {
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
}

func (tp *TermParser) GetCurBuf() string {
	return cur_buffer
}

func (tp *TermParser) GetPrefixString() string {
	return tp.prefix_string
}

func (tp *TermParser) Run() {
	defer handleExit()
	tp.p.Run()
}

func makeParser() TermParser {
	var tp TermParser

	tp.prefix_string = ">>> "
	tp.p = prompt.New(
		tp.executor,
		tp.completer,
		prompt.OptionSetExitCheckerOnInput(tp.exitChecker),
		prompt.OptionPrefix(tp.prefix_string),
		prompt.OptionPrefixTextColor(prompt.DarkGreen),
		prompt.OptionPreviewSuggestionTextColor(prompt.White),
		prompt.OptionSelectedSuggestionTextColor(prompt.DarkGreen),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionDescriptionBGColor(prompt.White),
	)

	return tp
}
