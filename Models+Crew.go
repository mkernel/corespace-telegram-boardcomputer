package main

func (crew crew) isCurrent() bool {
	return crew.ID == activeCrewID
}
