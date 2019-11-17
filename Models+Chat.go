package main

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (this_chat chat) hasUnread() bool {
	var unread uint
	database.Model(&message{}).Where("chat_id = ? and read = ?", this_chat.ID, false).Count(&unread)
	return unread > 0
}

func (this_chat chat) isLinked() bool {
	var count uint
	database.Model(&crew{}).Where("chat_id = ?", this_chat.ID).Count(&count)
	return count > 0
}

func (this_chat chat) FetchCrew() crew {
	var foundCrew crew
	database.Model(&crew{}).Where("chat_id = ?", this_chat.ID).First(&foundCrew)
	return foundCrew
}

func (this_chat chat) isCurrent() bool {
	return this_chat.ID == activeChatID
}

func (this_chat chat) fetchMessages() []message {
	var messages []message
	database.Model(&message{}).Where(&message{ChatID: this_chat.ID}).Order("id asc").Find(&messages)
	return messages
}

func (this_chat chat) sendMessage(text string) {
	msg := tgbotapi.NewMessage(this_chat.TelegramChatID, text)
	tgbot.Send(msg)
	newmsg := message{Inbound: false, ChatID: this_chat.ID, Text: text, Read: true, Date: int(time.Now().Unix())}
	database.Create(&newmsg)
	if activeChatID == this_chat.ID {
		output <- newmsg.toString()
	}
	if this_chat.Admin == false {
		var chats []chat
		database.Where(&chat{Admin: true}).Find(&chats)
		for _, admin := range chats {
			admin.sendMessage(fmt.Sprintf("[O] <%s> %s", this_chat.TelegramUserName, text))
		}
	}
}
