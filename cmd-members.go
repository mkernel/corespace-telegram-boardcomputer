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
	output(func(print printer) {
		print("members ls - lists all members")
		print("members rm ID - removes a member")
		print("members new - creates a new member")
		print("members update ID - updates the image of a member")

	})
}

func (command membersCmd) Execute(args []string) error {
	if activeCrewID == 0 {
		output(func(print printer) { print("no crew selected") })
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
	output(func(print printer) {
		for _, member := range members {
			print(fmt.Sprintf("%d: %s", member.ID, member.Name))
		}
	})
	return nil
}

func (membersCmd) rm(args []string) error {
	var member member
	var id int
	id, _ = strconv.Atoi(args[0])
	database.First(&member, uint(id))
	database.Delete(&member)
	output(func(print printer) {
		print(fmt.Sprintf("%s deleted.", member.Name))
	})
	return nil
}

func (command membersCmd) new(args []string) error {
	command.Mode = "new"
	command.Name = ""
	var casted cmdlinesink = command
	inputfocus = &casted
	output(func(print printer) {
		print("Name?")
	})
	return nil
}

func (command membersCmd) update(args []string) error {
	command.Mode = "update"
	id, _ := strconv.Atoi(args[0])
	command.MemberID = uint(id)
	var casted cmdlinesink = command
	inputfocus = &casted
	output(func(print printer) {
		print("Full path to new Character Board (PNG only)?")
	})
	return nil
}

func (command membersCmd) TextEntered(data string) error {
	if command.Mode == "new" {
		if command.Name == "" {
			command.Name = data
			output(func(print printer) {
				print(data)
				print("Full path to Character Board (PNG only)?")
			})
			var casted cmdlinesink = command
			inputfocus = &casted
		} else {
			newMember := member{CrewID: activeCrewID, Name: command.Name}
			database.Save(&newMember)
			filename := newMember.Filename()
			source, _ := os.Open(data)
			defer source.Close()
			destination, _ := os.Create(filename)
			defer destination.Close()
			io.Copy(destination, source)
			output(func(print printer) {
				print(data)
				print("Member created")
			})
		}
	} else if command.Mode == "update" {
		var member member
		database.First(&member, command.MemberID)
		filename := member.Filename()
		os.Remove(filename)
		source, _ := os.Open(data)
		defer source.Close()
		destination, _ := os.Create(filename)
		defer destination.Close()
		io.Copy(destination, source)
		output(func(print printer) {
			print(data)
			print("Member updated")
		})
	}
	return nil
}
