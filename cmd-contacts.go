package main

import "strconv"

type contactsCmd struct {
	Mode      string
	Name      string
	ContactID uint
}

func (contactsCmd) Command() string {
	return "contacts"
}

func (contactsCmd) Description() string {
	return "manages contacts of the current crew"
}

func (contactsCmd) Help(args []string) {
	output(func(printer printer) {
		printer("contacts select NAME - selects the current contact")
		printer("contacts new NAME - creates a new contact")
		printer("contacts link NAME CREWID - creates a new contact as a proxy to another crew")
		printer("contacts rm NAME - removes a contact")
		printer("contacts update NAME - updates the description of a contact")
		printer("contacts info NAME - lists description of contact")
	})
}

func (cmd contactsCmd) Execute(args []string) error {
	if activeCrewID == 0 {
		return nil
	}
	if args[0] == "select" {
		return cmd.cmdSelect(args[1:])
	} else if args[0] == "new" {
		return cmd.cmdNew(args[1:])
	} else if args[0] == "link" {
		return cmd.cmdLink(args[1:])
	} else if args[0] == "rm" {
		return cmd.cmdRm(args[1:])
	} else if args[0] == "update" {
		return cmd.cmdUpdate(args[1:])
	} else if args[0] == "info" {
		return cmd.cmdInfo(args[1:])
	}
	return nil
}

func (contactsCmd) cmdSelect(args []string) error {
	if args[0] == "_" {
		activeContactID = 0
		updateSidebar()
		return nil
	}
	var hit contact
	filter := contact{OwnerID: activeCrewID, Name: args[0]}
	database.Where(&filter).First(&hit)
	activeChatID = 0
	activeContactID = hit.ID
	updateSidebar()
	outputView.SetOrigin(0, 0)
	outputView.SetCursor(0, 0)
	outputView.Clear()
	mails := hit.fetchSpacemail()
	for _, mail := range mails {
		mail.print(outputView)
	}
	return nil
}

func (cmd contactsCmd) cmdNew(args []string) error {
	cmd.Mode = "new"
	cmd.Name = args[0]
	var casted cmdlinesink = cmd
	inputfocus = &casted
	output(func(printer printer) {
		printer("Description?")
	})
	return nil
}

func (contactsCmd) cmdLink(args []string) error {
	destinationid, _ := strconv.Atoi(args[1])
	newContact := contact{OwnerID: activeCrewID, Name: args[0], CrewID: uint(destinationid)}
	database.Create(&newContact)
	updateSidebar()
	return nil
}

func (contactsCmd) cmdRm(args []string) error {
	var hit contact
	filter := contact{OwnerID: activeCrewID, Name: args[0]}
	database.Where(&filter).First(&hit)
	database.Delete(&hit)
	updateSidebar()
	return nil
}

func (cmd contactsCmd) cmdUpdate(args []string) error {
	var hit contact
	filter := contact{OwnerID: activeCrewID, Name: args[0]}
	database.Where(&filter).First(&hit)
	cmd.Mode = "update"
	cmd.ContactID = hit.ID
	var casted cmdlinesink = cmd
	inputfocus = &casted
	output(func(printer printer) {
		printer("Description?")
	})
	return nil
}

func (contactsCmd) cmdInfo(args []string) error {
	var hit contact
	filter := contact{OwnerID: activeCrewID, Name: args[0]}
	database.Where(&filter).First(&hit)
	output(func(printer printer) {
		printer(hit.Description)
	})
	return nil
}

func (cmd contactsCmd) TextEntered(data string) error {
	if cmd.Mode == "new" {
		output(func(printer printer) {
			printer(data)
		})
		contact := contact{OwnerID: activeCrewID, Name: cmd.Name, Description: data}
		database.Create(&contact)
		updateSidebar()
	} else if cmd.Mode == "update" {
		output(func(printer printer) {
			printer(data)
		})
		var contact contact
		database.First(&contact, cmd.ContactID)
		contact.Description = data
		database.Save(&contact)
	}
	return nil
}
