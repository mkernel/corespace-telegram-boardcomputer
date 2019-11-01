package main

func (contact contact) isCurrent() bool {
	return contact.ID == activeContactID
}

func (contact contact) fetchSpacemail() []spacemail {
	var mails []spacemail
	filter := spacemail{CrewID: contact.CrewID, ContactID: contact.ID}
	database.Where(&filter).Order("date asc").Find(&mails)
	return mails
}

func (contact contact) numCrewUnread() int {
	var numUnread uint
	database.Model(&spacemail{}).Where("crew_id = ? and contact_id = ? and inbound = ? and read = ?", contact.CrewID, contact.ID, true, false).Count(&numUnread)
	return int(numUnread)
}

func (contact contact) numContactUnread() int {
	var numUnread uint
	database.Model(&spacemail{}).Where("crew_id = ? and contact_id = ? and inbound = ? and read = ?", contact.CrewID, contact.ID, false, false).Count(&numUnread)
	return int(numUnread)
}
