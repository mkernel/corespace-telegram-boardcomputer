package main

func (crew crew) isCurrent() bool {
	return crew.ID == activeCrewID
}

func (crew crew) fetchInventory() []item {
	var items []item
	filter := item{CrewID: crew.ID}
	database.Where(&filter).Find(&items)
	return items
}

func (crew crew) balance() float64 {
	var transactions []transaction
	filter := transaction{CrewID: crew.ID}
	database.Where(&filter).Find(&transactions)
	var balance float64
	for _, tx := range transactions {
		balance += tx.Value
	}
	return balance
}
func (crew crew) fetchContacts() []contact {
	var contacts []contact
	filter := contact{OwnerID: crew.ID}
	database.Where(&filter).Order("name asc").Find(&contacts)
	return contacts
}
