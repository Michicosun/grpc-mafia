package client

import "fmt"

var Printer = TermPrinter{}

type TermPrinter struct {
	// some state
}

func (tp *TermPrinter) PrintLine(from string, msg string) {
	fmt.Printf("\033[2K\r%s > %s\n", from, msg)
	fmt.Printf("%s%s", Parser.GetPrefixString(), Parser.GetCurBuf())
}

func (tp *TermPrinter) Init() {
	// init state
}
