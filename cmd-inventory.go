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
	output(func(print printer) {
		print("inventory ls - lists the inventory")
		print("inventory new - adds a new item to the inventory")
		print("inventory rm ID - removes an item")
		print("inventory update ID - updates the description of an item")
	})
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
	output(func(print printer) {
		for _, item := range items {
			print(fmt.Sprintf("%d: %s", item.ID, item.Name))
			print(item.Description)
		}
	})
	return nil
}

func (cmd inventoryCmd) new(args []string) error {
	cmd.Mode = "new"
	cmd.Name = ""
	var casted cmdlinesink = cmd
	inputfocus = &casted
	output(func(print printer) {
		print("Name?")
	})
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
	output(func(print printer) {
		print("New Description?")
	})
	return nil
}

func (cmd inventoryCmd) TextEntered(data string) error {
	if cmd.Mode == "new" {
		if cmd.Name == "" {
			cmd.Name = data
			output(func(print printer) {
				print(data)
				print("Description?")
			})
			var casted cmdlinesink = cmd
			inputfocus = &casted
		} else {
			output(func(print printer) {
				print(data)
			})
			newItem := item{CrewID: activeCrewID, Name: cmd.Name, Description: data}
			database.Save(&newItem)
		}
	} else if cmd.Mode == "update" {
		var item item
		database.First(&item, cmd.CurrentItem)
		item.Description = data
		output(func(print printer) {
			print(data)
		})
		database.Save(&item)
	}
	return nil
}
