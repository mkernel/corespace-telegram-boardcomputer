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
	conference := conference{InvolvedCrews: []crew{origin}, InvolvedNSCs: nscs, RingingCrews: crews}
	//TODO: set up ringing for everyone
	//TODO: generate join call for all NSCs. ?
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

func (mgr *conferenceMgr) hangup(crew crew) {
	conference := mgr.findConferenceForCrew(crew)
	if conference == nil {
		return
	}
	conference.hangup(crew)
	//TODO: check if there are crews left. if not, discard the conference
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
			return
		}
	}
}

func (cf *conference) hangup(crew crew) {
	for i, calling := range cf.InvolvedCrews {
		if calling.ID == crew.ID {
			cf.InvolvedCrews = append(cf.InvolvedCrews[:i], cf.InvolvedCrews[i+1:]...)
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
		//TODO: update the UI
	}
}

func (cf *conference) transmitForContact(contact contact, text string) {
	//TODO: we need to save the message for the contact himself.
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
			//TODO: update the UI
		}
	}
}
