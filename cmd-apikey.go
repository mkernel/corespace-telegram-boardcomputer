package main

import (
	"fmt"
)

type apiKeyCmd struct{}

func (apiKeyCmd) Command() string {
	return "apikey"
}

func (apiKeyCmd) Description() string {
	return "sets the telegram api key."
}

func (apiKeyCmd) Help(_ []string) {
	fmt.Fprintln(outputView, "No help available")
}

func (apiKeyCmd) Execute(args []string) error {
	var settings globalSettings
	database.First(&settings)
	if len(args) != 1 {
		if settings.APIKey == "" {
			fmt.Fprintln(outputView, "No API key set.")
		} else {
			//print the api key instead of setting it.
			fmt.Fprintf(outputView, "The API key is: %s\n", settings.APIKey)
		}
	} else {
		settings.APIKey = args[0]
		database.Save(&settings)
		fmt.Fprintln(outputView, "API key updated. Restart to take effect.")
	}
	return nil
}
