package main

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
	output <- "contacts select NAME - selects the current contact"
	output <- "contacts new NAME - creates a new contact"
	output <- "contacts link NAME CREWID - creates a new contact as a proxy to another crew"
	output <- "contacts rm NAME - removes a contact"
	output <- "contacts update NAME - updates the description of a contact"
	output <- "contacts info NAME - lists description of contact"
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
	filter := contact{CrewID: activeCrewID, Name: args[0]}
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
	output <- "Description?"
	return nil
}

func (contactsCmd) cmdLink(args []string) error {
	output <- "This has not been implemented yet. Will do that up next"
	//TODO: implement linking crews
	return nil
}

func (contactsCmd) cmdRm(args []string) error {
	var hit contact
	filter := contact{CrewID: activeCrewID, Name: args[0]}
	database.Where(&filter).First(&hit)
	database.Delete(&hit)
	updateSidebar()
	return nil
}

func (cmd contactsCmd) cmdUpdate(args []string) error {
	var hit contact
	filter := contact{CrewID: activeCrewID, Name: args[0]}
	database.Where(&filter).First(&hit)
	cmd.Mode = "udpate"
	cmd.ContactID = hit.ID
	var casted cmdlinesink = cmd
	inputfocus = &casted
	return nil
}

func (contactsCmd) cmdInfo(args []string) error {
	var hit contact
	filter := contact{CrewID: activeCrewID, Name: args[0]}
	database.Where(&filter).First(&hit)
	output <- hit.Description
	return nil
}

func (cmd contactsCmd) TextEntered(data string) error {
	if cmd.Mode == "new" {
		output <- data
		contact := contact{CrewID: activeCrewID, Name: cmd.Name, Description: data}
		database.Create(&contact)
		updateSidebar()
	}
	if cmd.Mode == "update" {
		output <- data
		var contact contact
		database.First(&contact, cmd.ContactID)
		contact.Description = data
		database.Save(&contact)
	}
	return nil
}
