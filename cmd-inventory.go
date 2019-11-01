package main

import (
	"fmt"
	"strconv"
)

type inventoryCmd struct {
	Mode        string
	Name        string
	CurrentItem uint
}

func (inventoryCmd) Command() string {
	return "inventory"
}

func (inventoryCmd) Description() string {
	return "allows managing the current crews inventory"
}

func (inventoryCmd) Help(args []string) {
	output <- "inventory ls - lists the inventory"
	output <- "inventory new - adds a new item to the inventory"
	output <- "inventory rm ID - removes an item"
	output <- "inventory update ID - updates the description of an item"
}

func (cmd inventoryCmd) Execute(args []string) error {
	if activeCrewID == 0 {
		return nil
	}
	if args[0] == "ls" {
		return cmd.ls(args[1:])
	} else if args[0] == "new" {
		return cmd.new(args[1:])
	} else if args[0] == "rm" {
		return cmd.rm(args[1:])
	} else if args[0] == "update" {
		return cmd.update(args[1:])
	}
	return nil
}

func (inventoryCmd) ls(args []string) error {
	var items []item
	filtered := item{CrewID: activeCrewID}
	database.Where(&filtered).Find(&items)
	for _, item := range items {
		output <- fmt.Sprintf("%d: %s", item.ID, item.Name)
		output <- item.Description
	}
	return nil
}

func (cmd inventoryCmd) new(args []string) error {
	cmd.Mode = "new"
	cmd.Name = ""
	var casted cmdlinesink = cmd
	inputfocus = &casted
	output <- "Name?"
	return nil
}

func (cmd inventoryCmd) rm(args []string) error {
	id, _ := strconv.Atoi(args[0])
	var item item
	database.First(&item, uint(id))
	database.Delete(&item)
	return nil
}

func (cmd inventoryCmd) update(args []string) error {
	id, _ := strconv.Atoi(args[0])
	cmd.Mode = "update"
	cmd.CurrentItem = uint(id)
	var casted cmdlinesink = cmd
	inputfocus = &casted
	output <- "New Description?"
	return nil
}

func (cmd inventoryCmd) TextEntered(data string) error {
	if cmd.Mode == "new" {
		if cmd.Name == "" {
			cmd.Name = data
			output <- data
			output <- "Description?"
			var casted cmdlinesink = cmd
			inputfocus = &casted
		} else {
			output <- data
			newItem := item{CrewID: activeCrewID, Name: cmd.Name, Description: data}
			database.Save(&newItem)
		}
	} else if cmd.Mode == "update" {
		var item item
		database.First(&item, cmd.CurrentItem)
		item.Description = data
		output <- data
		database.Save(&item)
	}
	return nil
}
