package main

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type exitCmd struct{}

func (exitCmd) Command() string {
	return "exit"
}

func (exitCmd) Description() string {
	return "Shuts down the boardcomputer management console"
}

func (exitCmd) Help(_ []string) {
	fmt.Fprintln(outputView, "No help available")
}

func (exitCmd) Execute(_ []string) error {
	return gocui.ErrQuit
}
