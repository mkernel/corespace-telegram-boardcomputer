package main

import (
	"fmt"
)

type ApiKeyCmd struct {}


func(_ ApiKeyCmd) Command() string {
	return "apikey"
}

func(_ ApiKeyCmd) Description() string {
	return "sets the telegram api key."
}

func(_ ApiKeyCmd) Help(_ []string) {
	fmt.Fprintln(OutputView,"No help available");
}

func(_ ApiKeyCmd) Execute(args []string) error {
	var settings GlobalSettings;
	Database.First(&settings);
	if(len(args) != 1) {
		if(settings.APIKey == "") {
			fmt.Fprintln(OutputView,"No API key set.");
		} else {
			//print the api key instead of setting it.
			fmt.Fprintf(OutputView,"The API key is: %s\n",settings.APIKey);
		}
	} else {
		settings.APIKey = args[0];
		Database.Save(&settings);
		fmt.Fprintln(OutputView,"API key updated. Restart to take effect.");
	}
	return nil;
}
