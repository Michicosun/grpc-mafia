package main

import (
	"fmt"
	"os"
	"os/exec"

	prompt "github.com/c-bata/go-prompt"
)

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "users", Description: "Store the username and age"},
		{Text: "articles", Description: "Store the article text posted by user"},
		{Text: "comments", Description: "Store the text commented to articles"},
		{Text: "groups", Description: "Combine users with specific rules"},
		{Text: "exit", Description: "quite the program"},
	}
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func handleExit() {
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
}

func main() {
	defer handleExit()

	for {
		in := prompt.Input(">>> ", completer,
			prompt.OptionPreviewSuggestionTextColor(prompt.White),
			prompt.OptionPrefixTextColor(prompt.DarkGreen),
			prompt.OptionSelectedSuggestionTextColor(prompt.DarkGreen),
			prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
			prompt.OptionDescriptionBGColor(prompt.White),
		)

		fmt.Println("Your input: " + in)

		if in == "exit" {
			break
		}
	}
}
