package main

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (chat chat) hasUnread() bool {
	var unread uint
	database.Model(&message{}).Where("chat_id = ? and read = ?", chat.ID, false).Count(&unread)
	return unread > 0
}

func (chat chat) isLinked() bool {
	var count uint
	database.Model(&crew{}).Where("chat_id = ?", chat.ID).Count(&count)
	return count > 0
}

func (chat chat) FetchCrew() crew {
	var foundCrew crew
	database.Model(&crew{}).Where("chat_id = ?", chat.ID).First(&foundCrew)
	return foundCrew
}

func (chat chat) isCurrent() bool {
	return chat.ID == activeChatID
}

func (chat chat) fetchMessages() []message {
	var messages []message
	database.Model(&message{}).Where(&message{ChatID: chat.ID}).Order("id asc").Find(&messages)
	return messages
}

func (chat chat) sendMessage(text string) {
	msg := tgbotapi.NewMessage(chat.TelegramChatID, text)
	tgbot.Send(msg)
	newmsg := message{Inbound: false, ChatID: chat.ID, Text: text, Read: true, Date: int(time.Now().Unix())}
	database.Create(&newmsg)
	if activeChatID == chat.ID {
		output <- newmsg.toString()
	}

}
