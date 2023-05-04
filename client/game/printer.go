package game

import (
	"fmt"
)

func PrintLine(from string, msg string, interactor IInteractor) {
	fmt.Printf("\033[2K\r%s > %s\n", from, msg)
	RefreshLine(interactor)
}

func RefreshLine(interactor IInteractor) {
	fmt.Printf("\033[2K\r%s%s\r", interactor.GetPrefixString(), interactor.GetCurBuf())
}
