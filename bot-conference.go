package main

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
	//TODO: implement /call
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
	worker.Commands = []botCommand{
		botHelpCmd{},
		botHangupCmd{},
	}
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
	conferences.hangup(worker.Chat.FetchCrew())
}
