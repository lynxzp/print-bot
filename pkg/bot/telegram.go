package bot

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lynxzp/print-bot/pkg/bot/text"
	"log"
)

type Config struct {
	Token       string
	AdminChatId int64
}

var config Config

func Run(cfg Config) {
	config = cfg
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			process1NewMessage(bot, update.Message)
			continue
		}

		if update.CallbackQuery != nil {
			process1Callback(bot, update.CallbackQuery)
			continue
		}
		j, _ := json.Marshal(update)
		log.Println("WW Received update without message and callback, update=", string(j))
	}
}

func process1NewMessage(bot *tgbotapi.BotAPI, m *tgbotapi.Message) {
	lang := text.SelectLang(m.From.LanguageCode)
	if len(m.Photo) != 0 {
		process2Photo(bot, m, lang)
		return
	}
	if m.Document != nil {
		process2Document(bot, m, lang)
		return
	}
	if m.Text == "/start" {
		process2StartMessage(bot, m, lang)
		return
	}
	process2TextMessage(bot, m, lang)
}

func process1Callback(bot *tgbotapi.BotAPI, c *tgbotapi.CallbackQuery) {
	lang := text.SelectLang(c.From.LanguageCode)
	if c.Data == "to support" {
		process2ForwardToSupport(bot, c, lang)
	}
	log.Println("From " + c.From.UserName + " received callback. Data: " + c.Data)
}

func process2ForwardToSupport(bot *tgbotapi.BotAPI, c *tgbotapi.CallbackQuery, lang string) {
	msg := tgbotapi.NewForward(config.AdminChatId, c.From.ID, c.Message.ReplyToMessage.MessageID)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("failed to forward message to support:", err)
	}
	edit := tgbotapi.NewEditMessageText(c.From.ID, c.Message.MessageID, text.WasSentToSupport[lang])
	_, err = bot.Send(edit)
	if err != nil {
		log.Println("failed to edit message:", err)
	}
}

func process2TextMessage(bot *tgbotapi.BotAPI, m *tgbotapi.Message, lang string) {
	msg := tgbotapi.NewMessage(m.Chat.ID, text.ReplyUserSentText[lang])
	msg.ReplyToMessageID = m.MessageID
	b := tgbotapi.NewInlineKeyboardButtonData(text.SendToSupport[lang], "to support")
	markup := tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{b})
	msg.ReplyMarkup = markup
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("failed to send message:", err)
	}
}

func process2StartMessage(bot *tgbotapi.BotAPI, m *tgbotapi.Message, lang string) {
	msg := tgbotapi.NewMessage(m.Chat.ID, text.ReplyStartMessage[lang])
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("failed to send message:", err)
	}
}

func process2Photo(bot *tgbotapi.BotAPI, m *tgbotapi.Message, lang string) {
	msg := tgbotapi.NewMessage(m.Chat.ID, "")
	log.Println("From " + m.From.UserName + " received photo")
	msg.Text = "From " + m.From.UserName + " received photo"
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("failed to send message:", err)
	}
}

func process2Document(bot *tgbotapi.BotAPI, m *tgbotapi.Message, lang string) {
	msg := tgbotapi.NewMessage(m.Chat.ID, "")
	log.Println("From " + m.From.UserName + " received document")
	msg.Text = "From " + m.From.UserName + " received document"
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("failed to send message:", err)
	}
}
