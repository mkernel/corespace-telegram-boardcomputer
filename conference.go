package main

import (
	"fmt"
	"time"
)

type conference struct {
	InvolvedCrews []crew
	InvolvedNSCs  []contact
	RingingCrews  []crew
}

type conferenceMgr struct {
	Conferences []conference
}

var conferences conferenceMgr

func setupConferenceMgr() {
	conferences.setup()
}

func (mgr *conferenceMgr) setup() {
	mgr.Conferences = make([]conference, 0)
}

func (mgr *conferenceMgr) call(origin crew, crews []crew, nscs []contact) {
	originworker := getWorker(origin.ChatID)
	originworker.Chat.OpenConnection = true
	database.Save(&originworker.Chat)
	conference := conference{InvolvedCrews: []crew{origin}, InvolvedNSCs: nscs, RingingCrews: crews}
	for _, crew := range conference.RingingCrews {
		worker := getWorker(crew.ChatID)
		worker.Chat.OpenConnection = true
		database.Save(&worker.Chat)
		worker.setCommandSet([]botCommand{
			botHelpCmd{},
			botAcceptCmd{},
			botRejectCmd{},
		})
		worker.Chat.sendMessage(fmt.Sprintf("Eingehende Konferenzanfrage von %s. Mit /accept annehmen, mit /reject ablehnen.", origin.Name))

	}
	for _, contact := range nscs {
		conference.transmitForContact(contact, "*SYSTEM* Verbindung hergestellt.")
	}
	mgr.Conferences = append(mgr.Conferences, conference)
}

func (mgr *conferenceMgr) findConferenceForCrew(crew crew) *conference {
	for i, conference := range mgr.Conferences {
		for _, participatingcrew := range conference.InvolvedCrews {
			if participatingcrew.ID == crew.ID {
				return &mgr.Conferences[i]
			}
		}
		for _, ringingcrew := range conference.RingingCrews {
			if ringingcrew.ID == crew.ID {
				return &mgr.Conferences[i]
			}
		}
	}
	return nil
}

func (mgr *conferenceMgr) findConferenceForContact(contact contact) *conference {
	for i, conference := range mgr.Conferences {
		for _, oncall := range conference.InvolvedNSCs {
			if oncall.ID == contact.ID {
				return &mgr.Conferences[i]
			}
		}
	}
	return nil
}

func (mgr *conferenceMgr) acceptCall(crew crew) {
	conference := mgr.findConferenceForCrew(crew)
	if conference == nil {
		return
	}
	conference.accept(crew)
}

func (mgr *conferenceMgr) rejectCall(crew crew) {
	conference := mgr.findConferenceForCrew(crew)
	if conference == nil {
		return
	}
	conference.reject(crew)
}

func (mgr *conferenceMgr) hangup(crew crew) {
	conference := mgr.findConferenceForCrew(crew)
	if conference == nil {
		return
	}
	conference.hangup(crew)
	if len(conference.InvolvedCrews) == 0 {
		for _, crew := range conference.RingingCrews {
			conference.reject(crew)
		}
		for i, conf := range mgr.Conferences {
			if len(conf.InvolvedCrews) == 0 && len(conf.RingingCrews) == 0 {
				mgr.Conferences = append(mgr.Conferences[:i], mgr.Conferences[i+1:]...)
			}
		}
	}
}

func (mgr *conferenceMgr) isCrewInOngoingCall(crew crew) bool {
	return mgr.findConferenceForCrew(crew) != nil
}

func (mgr *conferenceMgr) isContactInOngoingCall(contact contact) bool {
	return mgr.findConferenceForContact(contact) != nil
}

func (mgr *conferenceMgr) transmitFromCrew(crew crew, text string) {
	conference := mgr.findConferenceForCrew(crew)
	if conference != nil {
		conference.transmitForCrew(crew, text)
	}
}

func (mgr *conferenceMgr) transmitFromContact(contact contact, text string) {
	conference := mgr.findConferenceForContact(contact)
	if conference != nil {
		conference.transmitForContact(contact, text)
	}
}

// conference API

func (cf *conference) accept(crew crew) {
	for i, ringing := range cf.RingingCrews {
		if ringing.ID == crew.ID {
			cf.RingingCrews = append(cf.RingingCrews[:i], cf.RingingCrews[i+1:]...)
			cf.InvolvedCrews = append(cf.InvolvedCrews, crew)
			worker := getWorker(crew.ChatID)
			worker.setCommandSet([]botCommand{
				botHelpCmd{},
				botHangupCmd{},
			})
			return
		}
	}
}

func (cf *conference) reject(crew crew) {
	for i, calling := range cf.RingingCrews {
		if calling.ID == crew.ID {
			cf.transmitForCrew(crew, "*SYSTEM* Verbindungsaufbau abgelehnt.")
			cf.RingingCrews = append(cf.RingingCrews[:i], cf.RingingCrews[i+1:]...)
			worker := getWorker(crew.ChatID)
			worker.setDefaultCommandSet()
			return
		}
	}
}

func (cf *conference) hangup(crew crew) {
	for i, calling := range cf.InvolvedCrews {
		if calling.ID == crew.ID {
			cf.InvolvedCrews = append(cf.InvolvedCrews[:i], cf.InvolvedCrews[i+1:]...)
			cf.transmitForCrew(crew, "*SYSTEM* Verbindung beendet.")
			worker := getWorker(crew.ChatID)
			worker.setDefaultCommandSet()
			return
		}
	}
}

func (cf *conference) transmitForCrew(crew crew, text string) {
	for _, oncall := range cf.InvolvedCrews {
		if oncall.ID != crew.ID && oncall.ChatID != 0 {
			var chat chat
			database.First(&chat, oncall.ChatID)
			chat.sendMessage(fmt.Sprintf("<%s> %s", crew.Name, text))
		}
	}
	for _, contact := range cf.InvolvedNSCs {
		spacemail := spacemail{CrewID: crew.ID, ContactID: contact.ID, Text: text, Inbound: false, Read: false, Date: int(time.Now().Unix())}
		database.Create(&spacemail)
		updateSidebar()
		if contact.ID == activeContactID {
			output <- spacemail.toString()
		}
	}
}

func (cf *conference) transmitForContact(contact contact, text string) {
	protocol := spacemail{CrewID: contact.OwnerID, ContactID: contact.ID, Text: text, Date: int(time.Now().Unix()), Inbound: true, Read: true}
	database.Create(&protocol)

	for _, oncall := range cf.InvolvedCrews {
		if oncall.ChatID != 0 {
			var chat chat
			database.First(&chat, oncall.ChatID)
			chat.sendMessage(fmt.Sprintf("<%s> %s", contact.Name, text))
		}
	}
	for _, oncall := range cf.InvolvedNSCs {
		if oncall.ID != contact.ID {
			spacemail := spacemail{CrewID: oncall.OwnerID, ContactID: oncall.ID, Text: fmt.Sprintf("<%s> %s", contact.Name, text), Inbound: false, Read: false, Date: int(time.Now().Unix())}
			database.Create(&spacemail)
			updateSidebar()
			if contact.ID == activeContactID {
				output <- spacemail.toString()
			}
		}
	}
}
