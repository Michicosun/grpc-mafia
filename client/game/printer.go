package game

import (
	"fmt"
)

func PrintLine(from string, msg string, interactor IInteractor) {
	fmt.Printf("\033[2K\r%s > %s\n", from, msg)
	fmt.Printf("%s%s", interactor.GetPrefixString(), interactor.GetCurBuf())
}
