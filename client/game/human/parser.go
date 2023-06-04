package game

import (
	"fmt"
	"grpc-mafia/client/game"
	"os"
	"os/exec"

	"github.com/c-bata/go-prompt"
)

type humanInteractor struct {
	login      string
	prefix     string
	cur_buffer string
	p          *prompt.Prompt
}

func handleExit() {
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
}

func (hi *humanInteractor) SetLogin(login string) {
	hi.login = login
}

func (hi *humanInteractor) GetCurBuf() string {
	return hi.cur_buffer
}

func (hi *humanInteractor) GetPrefixString() string {
	green := "\033[32m"
	reset := "\033[0m"
	return fmt.Sprintf("%s%s%s", green, hi.prefix, reset)
}

func (hi *humanInteractor) Signal() {
	if game.Session.GetState() == game.Undefined {
		hi.prefix = ">>> "
	} else {
		hi.prefix = fmt.Sprintf("[%s]-[%s] >>> ", game.Session.Role.String(), getAliveString())
	}
}

func (hi *humanInteractor) Run() {
	defer handleExit()
	hi.p.Run()
}

func (hi *humanInteractor) changeLivePrefix() (string, bool) {
	return hi.prefix, true
}

func MakeHumanInteractor() *humanInteractor {
	hi := humanInteractor{}

	hi.cur_buffer = ""
	hi.prefix = ">>> "
	hi.p = prompt.New(
		hi.Executor,
		hi.Completer,
		prompt.OptionSetExitCheckerOnInput(hi.exitChecker),
		prompt.OptionLivePrefix(hi.changeLivePrefix),
		prompt.OptionPrefixTextColor(prompt.DarkGreen),
		prompt.OptionPreviewSuggestionTextColor(prompt.White),
		prompt.OptionSelectedSuggestionTextColor(prompt.DarkGreen),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionDescriptionBGColor(prompt.White),
	)

	return &hi
}

func getAliveString() string {
	if game.Session.GetState() == game.Ghost {
		return "Ghost"
	}
	return "Alive"
}
