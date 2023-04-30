package client

import "fmt"

var Printer = TermPrinter{}

type TermPrinter struct {
	// some state i guess
}

func (tp *TermPrinter) PrintLine(line string) {
	fmt.Printf("\033[2K\r%s\n", line)
	fmt.Printf("%s%s", Parser.GetPrefixString(), Parser.GetCurBuf())
}
