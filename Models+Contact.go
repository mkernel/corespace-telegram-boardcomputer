package main

import "time"

func (me contact) isCurrent() bool {
	return me.ID == activeContactID
}

func (me contact) fetchSpacemail() []spacemail {
	var mails []spacemail
	filter := spacemail{ContactID: me.ID}
	database.Where(&filter).Order("date asc").Find(&mails)
	return mails
}

func (me contact) numContactUnread() int {
	var numUnread uint
	database.Model(&spacemail{}).Where("contact_id = ? and inbound = ? and read = ?", me.ID, false, false).Count(&numUnread)
	return int(numUnread)
}

func (me contact) sendMessageToContact(text string) {
	mail := spacemail{CrewID: me.OwnerID, ContactID: me.ID, Inbound: false, Read: false, Date: int(time.Now().Unix()), Text: text}
	database.Create(&mail)
	if me.CrewID != 0 {
		var mirrorcontact contact
		database.Where(&contact{OwnerID: me.CrewID, CrewID: me.OwnerID}).First(&mirrorcontact)
		mirrorcontact.sendMessageToCrew(text)
	}
	if me.isCurrent() {
		output <- mail.toString()
	}
}

func (me contact) sendMessageToCrew(text string) {
	mirrormail := spacemail{CrewID: me.OwnerID, ContactID: me.ID, Inbound: true, Read: false, Date: int(time.Now().Unix()), Text: text}
	database.Create(&mirrormail)
	if me.isCurrent() {
		output <- mirrormail.toString()
	}
	var crew crew
	database.Preload("Chat").First(&crew, me.OwnerID)
	if crew.ChatID != 0 {
		crew.Chat.sendMessage(mirrormail.toString())
	}

}
