package main

type botCommand interface {
	Command() string
	Description() string
	Execute(*automationworker, []string)
}

type botDataSink interface {
	OnMessage(*automationworker, message)
}
