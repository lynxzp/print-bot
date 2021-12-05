package bot

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lynxzp/print-bot/pkg/bot/text"
	"github.com/lynxzp/print-bot/pkg/osenv"
	"log"
	"os"
	"strings"
)

type Config struct {
	Token       string
	AdminChatId int64
}

var (
	config Config
	bot    *tgbotapi.BotAPI
)

func Run(cfg Config) {
	config = cfg
	var err error
	bot, err = tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account @%s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			process1NewMessage(update.Message)
			continue
		}

		if update.CallbackQuery != nil {
			process1Callback(update.CallbackQuery)
			continue
		}
		j, _ := json.Marshal(update)
		log.Println("WW Received update without message and callback, update=", string(j))
	}
}

func process1NewMessage(m *tgbotapi.Message) {
	lang := text.SelectLang(m.From.LanguageCode)
	if len(m.Photo) != 0 {
		process2Photo(m, lang)
		return
	}
	if m.Document != nil {
		process2Document(m, lang)
		return
	}
	if m.Text == "/start" {
		process2StartMessage(m, lang)
		return
	}
	process2TextMessage(m, lang)
}

func process1Callback(c *tgbotapi.CallbackQuery) {
	lang := text.SelectLang(c.From.LanguageCode)
	if c.Data == "to support" {
		process2ForwardToSupport(c, lang)
	}
	log.Println("WW From " + getUserLink(c.From) + " received unprocessed callback. Data: " + c.Data)
}

func process2ForwardToSupport(c *tgbotapi.CallbackQuery, lang string) {
	msg := tgbotapi.NewForward(config.AdminChatId, c.From.ID, c.Message.ReplyToMessage.MessageID)
	sendMessage(msg)
	edit := tgbotapi.NewEditMessageText(c.From.ID, c.Message.MessageID, text.WasSentToSupport[lang])
	sendMessage(edit)
}

func process2TextMessage(m *tgbotapi.Message, lang string) {
	msg := tgbotapi.NewMessage(m.Chat.ID, text.ReplyUserSentText[lang])
	msg.ReplyToMessageID = m.MessageID
	b := tgbotapi.NewInlineKeyboardButtonData(text.SendToSupport[lang], "to support")
	markup := tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{b})
	msg.ReplyMarkup = markup
	sendMessage(msg)
}

func process2StartMessage(m *tgbotapi.Message, lang string) {
	msg := tgbotapi.NewMessage(m.Chat.ID, text.ReplyStartMessage[lang])
	sendMessage(msg)
}

func process2Photo(m *tgbotapi.Message, lang string) {
	log.Println("WW From " + getUserLink(m.From) + " received photo. Event unprocessed")
}

func process3PhotoAsDoc(m *tgbotapi.Message, lang string) {
	log.Println("WW From " + getUserLink(m.From) + " received photo as doc. Event unprocessed")
}

func process2Document(m *tgbotapi.Message, lang string) {
	mimes := strings.Split(m.Document.MimeType, "/")
	var mime0, _ string
	if len(mimes) == 2 {
		mime0 = mimes[0]
		_ = mimes[1]
	}
	if mime0 == "image" {
		process3PhotoAsDoc(m, lang)
		return
	}

	filenameExt := strings.Split(m.Document.FileName, ".")
	var ext string
	if len(filenameExt) > 1 {
		ext = filenameExt[len(filenameExt)-1]
	}
	if osenv.InConvertibleToPDF(ext) {
		process3DirectToPDF(m, lang)
		return
	}

	log.Println(m.Document.FileName, m.Document.MimeType)
	log.Println("WW From " + getUserLink(m.From) + " received document. Event unprocessed")
}

func process3DirectToPDF(m *tgbotapi.Message, lang string) {
	msg := tgbotapi.NewMessage(m.Chat.ID, text.ReplyUserSentText[lang])
	msg.ReplyToMessageID = m.MessageID
	msg.ReplyMarkup = makeDocumentKeyboard("", lang)
	sendMessage(msg)
}

func process4DirectToPDF(m *tgbotapi.Message, lang string) {

	// download
	url, err := bot.GetFileDirectURL(m.Document.FileID)
	if err != nil {
		log.Println("WW can't get direct file url fileId=", m.Document.FileID, " from ", getUserLink(m.From))
		return
	}
	path := osenv.GenFilePath(m.Document.FileUniqueID, "3", m.Document.FileName)
	err = osenv.DownloadFile(path, url)
	if err != nil {
		log.Println("WW can't download file ", url)
	}
	//goland:noinspection GoUnhandledErrorResult`
	defer os.Remove(path)
	log.Println(path, err)

	// convert
	//cmd := exec.Command(`c:\Program Files\PDFCreator\PDFCreator-cli.exe`, "PrintFile", "/File=")
}
