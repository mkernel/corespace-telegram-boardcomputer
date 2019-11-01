package main

import "fmt"

func (member member) Filename() string {
	return fmt.Sprintf("./assets/member_%d.png", member.ID)
}
