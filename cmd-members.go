package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

type membersCmd struct {
	Mode     string
	Name     string
	MemberID uint
}

func (membersCmd) Command() string {
	return "members"
}

func (membersCmd) Description() string {
	return "manages the crew members of the selected crew"
}

func (membersCmd) Help(args []string) {
	output <- "members ls - lists all members"
	output <- "members rm ID - removes a member"
	output <- "members new - creates a new member"
	output <- "members update ID - updates the image of a member"
}

func (command membersCmd) Execute(args []string) error {
	if activeCrewID == 0 {
		output <- "no crew selected"
	}
	if args[0] == "ls" {
		return command.ls(args[1:])
	} else if args[0] == "rm" {
		return command.rm(args[1:])
	} else if args[0] == "new" {
		return command.new(args[1:])
	} else if args[0] == "update" {
		return command.update(args[1:])
	}
	return nil
}

func (membersCmd) ls(args []string) error {
	var members []member
	filter := member{CrewID: activeCrewID}
	database.Where(&filter).Find(&members)
	for _, member := range members {
		output <- fmt.Sprintf("%d: %s", member.ID, member.Name)
	}
	return nil
}

func (membersCmd) rm(args []string) error {
	var member member
	var id int
	id, _ = strconv.Atoi(args[0])
	database.First(&member, uint(id))
	database.Delete(&member)
	output <- fmt.Sprintf("%s deleted.", member.Name)
	return nil
}

func (command membersCmd) new(args []string) error {
	command.Mode = "new"
	command.Name = ""
	var casted cmdlinesink = command
	inputfocus = &casted
	output <- "Name?"
	return nil
}

func (command membersCmd) update(args []string) error {
	command.Mode = "update"
	id, _ := strconv.Atoi(args[0])
	command.MemberID = uint(id)
	var casted cmdlinesink = command
	inputfocus = &casted
	output <- "Full path to new Character Board (PNG only)?"
	return nil
}

func (command membersCmd) TextEntered(data string) error {
	if command.Mode == "new" {
		if command.Name == "" {
			command.Name = data
			output <- data
			output <- "Full path to Character Board (PNG only)?"
			var casted cmdlinesink = command
			inputfocus = &casted
		} else {
			output <- data
			newMember := member{CrewID: activeCrewID, Name: command.Name}
			database.Save(&newMember)
			filename := fmt.Sprintf("assets/member_%d.png", newMember.ID)
			source, _ := os.Open(data)
			defer source.Close()
			destination, _ := os.Create(filename)
			defer destination.Close()
			io.Copy(destination, source)
			output <- "Member created"
		}
	} else if command.Mode == "update" {
		output <- data
		var member member
		database.First(&member, command.MemberID)
		filename := fmt.Sprintf("assets/member_%d.png", member.ID)
		os.Remove(filename)
		source, _ := os.Open(data)
		defer source.Close()
		destination, _ := os.Create(filename)
		defer destination.Close()
		io.Copy(destination, source)
		output <- "Member updated"
	}
	return nil
}
