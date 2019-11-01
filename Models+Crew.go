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
