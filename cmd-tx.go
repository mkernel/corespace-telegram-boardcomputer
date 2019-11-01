package main

import (
	"fmt"
	"strconv"
	"time"
)

type txCmd struct {
	Amount float64
}

func (txCmd) Command() string {
	return "tx"
}

func (txCmd) Description() string {
	return "Manages Transactions for the current crew"
}

func (txCmd) Help(args []string) {
	output <- "tx ls - lists all transactions and the balance"
	output <- "tx new VALUE - creates a new transaction"
}

func (cmd txCmd) Execute(args []string) error {
	if activeCrewID == 0 {
		return nil
	}
	if args[0] == "ls" {
		return cmd.ls(args[1:])
	} else if args[0] == "new" {
		return cmd.new(args[1:])
	}
	return nil
}

func (txCmd) ls(args []string) error {
	var transactions []transaction
	filter := transaction{CrewID: activeCrewID}
	database.Where(&filter).Order("date asc").Find(&transactions)
	var balance float64
	for _, tx := range transactions {
		output <- fmt.Sprintf("%.2f - %s", tx.Value, tx.Description)
		balance += tx.Value
	}
	output <- "---"
	output <- fmt.Sprintf("Balance: %.2f", balance)
	return nil
}

func (cmd txCmd) new(args []string) error {
	cmd.Amount, _ = strconv.ParseFloat(args[0], 64)
	var casted cmdlinesink = cmd
	inputfocus = &casted
	output <- "Description?"
	return nil
}

func (cmd txCmd) TextEntered(data string) error {
	output <- data
	tx := transaction{Date: int(time.Now().Unix()), CrewID: activeCrewID, Value: cmd.Amount, Description: data}
	database.Create(&tx)
	return nil
}
