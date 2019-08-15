package main;

import (
	"github.com/jroimartin/gocui"
	"fmt"
)

type ExitCmd struct {}

func(_ ExitCmd) Command() string {
	return "exit"
}

func(_ ExitCmd) Description() string {
	return "Shuts down the boardcomputer management console"
}

func(_ ExitCmd) Help(_ []string) {
	fmt.Fprintln(OutputView,"No help available");
}

func(_ ExitCmd) Execute(_ []string) error {
	return gocui.ErrQuit
}
