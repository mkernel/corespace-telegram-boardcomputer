package main

type botConnectCmd struct {
	ContactID uint
}

func (cmd botConnectCmd) Command() string {
	return "/connect"
}

func (cmd botConnectCmd) Description() string {
	return "Baut eine dauerhafte Verbindung zu einem Kontakt auf (Namen mit angeben)"
}

func (cmd botConnectCmd) Execute(worker *automationworker, args []string) {
	if len(args) == 0 {
		worker.Chat.sendMessage("Keinen Kontaktnamen angegeben")
		return
	}
	filter := contact{OwnerID: worker.Chat.FetchCrew().ID, Name: args[0]}
	var found contact
	database.Where(&filter).First(&found)
	if database.NewRecord(&found) {
		worker.Chat.sendMessage("Keinen passenden Kontakt in der Datenbank gefunden.")
		return
	}
	cmd.ContactID = found.ID
	worker.Chat.OpenConnection = true
	database.Save(&worker.Chat)
	worker.setCommandSet([]botCommand{
		botDisconnectCmd{},
		botHelpCmd{},
	})
	worker.Chat.sendMessage("Die Verbindung steht. Alles, was du schreibst, wird direkt übertragen. Benutze /disconnect um die Verbindung zu trennen. Das ist der einzige funktionierende Befehl.")
	var casted botDataSink = cmd
	worker.CurrentFocus = &casted
}

func (cmd botConnectCmd) OnMessage(worker *automationworker, msg message) {
	var ct contact
	database.First(&ct, cmd.ContactID)
	ct.sendMessageToContact(msg.Text)
	var casted botDataSink = cmd
	worker.CurrentFocus = &casted
}

type botDisconnectCmd struct{}

func (cmd botDisconnectCmd) Command() string {
	return "/disconnect"
}

func (cmd botDisconnectCmd) Description() string {
	return "Aktuelle Verbindung unterbrechen"
}

func (cmd botDisconnectCmd) Execute(worker *automationworker, artgs []string) {
	worker.CurrentFocus = nil
	worker.Chat.sendMessage("Verbindung getrennt")
	worker.setDefaultCommandSet()
}
