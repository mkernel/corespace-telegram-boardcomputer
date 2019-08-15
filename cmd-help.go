package main

import (
	"fmt"
)

type HelpCmd struct {}

func(_ HelpCmd) Command() string {
	return "help"
}

func(_ HelpCmd) Description() string {
	return "shows a list of available commands and executes their help"
}

func(_ HelpCmd) Help(_ []string) {
	fmt.Fprintln(OutputView,"Usage: help [Command]");
	fmt.Fprintln(OutputView,"without Command, lists all available commands");
}

func(_ HelpCmd) Execute(args []string) error {
	if(len(args)==0) {
		//this is the default case.
		for _,cmd := range List {
			fmt.Fprintf(OutputView,"%s - %s\n",cmd.Command(),cmd.Description());
		} 
	} else {
		command := args[0];
		for _,cmd := range List {
			if(cmd.Command() == command) {
				cmd.Help(args[1:]);
			}
		}
	}
	return nil;
}
