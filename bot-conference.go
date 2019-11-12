package main

import "fmt"

//basic idea: when you get one of those, the system updates the command list.

// /call implementation
type botCallCmd struct {
}

func (cmd botCallCmd) Command() string {
	return "/call"
}

func (cmd botCallCmd) Description() string {
	return "Baut eine direkte Kommunikationsverbindung auf. Einfach alle Kontaktnamen nacheinander angeben. (z.B. /call A B C)"
}

func (cmd botCallCmd) Execute(worker *automationworker, args []string) {
	var ringingCrews = make([]crew, 0)
	var ringingNSCs = make([]contact, 0)
	for _, name := range args {
		contact := fetchContactByName(name, worker.Chat.FetchCrew())
		if contact == nil {
			output <- fmt.Sprintf("%s ist nicht in der Kontaktliste", name)
			return
		}
		if contact.CrewID != 0 {
			//this is a linked contact.
			var crew crew
			database.First(&crew, contact.CrewID)
			ringingCrews = append(ringingCrews, crew)
		} else {
			ringingNSCs = append(ringingNSCs, *contact)
		}
	}
	//we have a list of contacts and crews.
	conferences.call(worker.Chat.FetchCrew(), ringingCrews, ringingNSCs)
	worker.Commands = []botCommand{
		botHelpCmd{},
		botHangupCmd{},
	}
	worker.Chat.sendMessage("Verbindung wird aufgebaut. Befehlsliste aktualisiert: mit /hangup Konferenz beenden.")
}

// /accept implementation
type botAcceptCmd struct {
}

func (cmd botAcceptCmd) Command() string {
	return "/accept"
}

func (cmd botAcceptCmd) Description() string {
	return "Nimmt eine eingehende Konferenzverbindung an."
}

func (cmd botAcceptCmd) Execute(worker *automationworker, args []string) {
	conferences.acceptCall(worker.Chat.FetchCrew())
	worker.Chat.sendMessage("Verbindung hergestellt. Alles geschriebene wird an alle Gesprächsteilnehmer übertragen. Mit /hangup kann die Verbindung getrennt werden.")
	conferences.transmitFromCrew(worker.Chat.FetchCrew(), "*SYSTEM* Verbindung hergestellt.")
}

// /reject implementation
type botRejectCmd struct {
}

func (cmd botRejectCmd) Command() string {
	return "/reject"
}

func (cmd botRejectCmd) Description() string {
	return "Lehnt eine eingehende Konferenzverbindung ab."
}

func (cmd botRejectCmd) Execute(worker *automationworker, args []string) {
	worker.setDefaultCommandSet()
	worker.Chat.sendMessage("Verbindung abgelehnt.")
	conferences.rejectCall(worker.Chat.FetchCrew())
}

// /hangup implementation
type botHangupCmd struct {
}

func (cmd botHangupCmd) Command() string {
	return "/hangup"
}

func (cmd botHangupCmd) Description() string {
	return "Beendet die Konferenzverbindung."
}

func (cmd botHangupCmd) Execute(worker *automationworker, args []string) {
	conferences.transmitFromCrew(worker.Chat.FetchCrew(), "*SYSTEM* Verbindung getrennt.")
	worker.Chat.sendMessage("Verbindung getrennt.")
	conferences.hangup(worker.Chat.FetchCrew())
}
