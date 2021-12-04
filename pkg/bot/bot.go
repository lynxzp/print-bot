package bot

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

func sendMessage(c tgbotapi.Chattable) {
	// todo: add more retries
	_, err := bot.Send(c)
	j, _ := json.Marshal(c)
	if err != nil {
		log.Println("EE failed to send message. Err:", err, ", msg:", j)
	}
}

func getUserLink(user *tgbotapi.User) string {
	return `{"Username":"` + user.UserName + `", "ID":` + strconv.FormatInt(user.ID, 10) + `}`
}
