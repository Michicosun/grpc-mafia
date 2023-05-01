package game

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/c-bata/go-prompt"
)

type humanInteractor struct {
	prefix_string string
	cur_buffer    string
	p             *prompt.Prompt
}

func handleExit() {
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
}

func (hi *humanInteractor) GetCurBuf() string {
	return hi.cur_buffer
}

func (hi *humanInteractor) GetPrefixString() string {
	green := "\033[32m"
	reset := "\033[0m"
	return fmt.Sprintf("%s%s%s", green, hi.prefix_string, reset)
}

func (hi *humanInteractor) Signal() {}

func (hi *humanInteractor) Run() {
	defer handleExit()
	hi.p.Run()
}

func MakeHumanInteractor() *humanInteractor {
	hi := humanInteractor{}

	hi.cur_buffer = ""
	hi.prefix_string = ">>> "
	hi.p = prompt.New(
		hi.Executor,
		hi.Completer,
		prompt.OptionSetExitCheckerOnInput(hi.exitChecker),
		prompt.OptionPrefix(hi.prefix_string),
		prompt.OptionPrefixTextColor(prompt.DarkGreen),
		prompt.OptionPreviewSuggestionTextColor(prompt.White),
		prompt.OptionSelectedSuggestionTextColor(prompt.DarkGreen),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionDescriptionBGColor(prompt.White),
	)

	return &hi
}
