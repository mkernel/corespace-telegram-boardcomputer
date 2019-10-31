package main

import (
	"fmt"
	"strconv"
)

type crewCmd struct {
	NewName string
	NewCode string
}

func (crewCmd) Command() string {
	return "crew"
}

func (crewCmd) Description() string {
	return "Manages crews"
}

func (crewCmd) Help(_ []string) {
	output <- "crew ls – list all crews"
	output <- "crew new - creates a new crew"
	output <- "crew rm ID - deletes a crew"
}

func (command crewCmd) Execute(args []string) error {
	if len(args) == 0 {
		command.Help(args)
	} else {
		if args[0] == "ls" {
			return command.ls(args[1:])
		}
		if args[0] == "new" {
			return command.new(args[1:])
		}
		if args[0] == "rm" {
			return command.rm(args[1:])
		}
	}
	return nil
}

func (crewCmd) ls(args []string) error {
	var crews []crew
	database.Find(&crews)
	for _, crew := range crews {
		output <- fmt.Sprintf("%d: %s (%s)", crew.ID, crew.Name, crew.Code)
		output <- crew.Description
	}
	return nil
}

func (crewCmd) rm(args []string) error {
	var crew crew
	id, _ := strconv.Atoi(args[0])
	database.First(&crew, id)
	database.Delete(&crew)
	output <- "Crew deleted."
	return nil
}

func (command crewCmd) new(args []string) error {
	var casted cmdlinesink = command
	inputfocus = &casted
	output <- "Crew Name?"
	return nil
}

func (command crewCmd) TextEntered(data string) error {
	if command.NewName == "" {
		command.NewName = data
		output <- data
		var casted cmdlinesink = command
		inputfocus = &casted
		output <- "Access Code?"
	} else if command.NewCode == "" {
		command.NewCode = data
		output <- data
		var casted cmdlinesink = command
		inputfocus = &casted
		output <- "Crew Description?"
	} else {
		output <- data
		crew := crew{Name: command.NewName, Description: data, Code: command.NewCode}
		database.Create(&crew)
		output <- fmt.Sprintf("Crew created with ID %d", crew.ID)
		command.NewCode = ""
		command.NewName = ""
	}
	return nil
}
