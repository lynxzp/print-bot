package text

var (
	ReplyStartMessage = map[string]string{
		"uk": "Привіт. Щоб роздрукувати документ чи фото просто відправ його мені.",
	}
	ReplyUserSentText = map[string]string{
		"uk": "Привіт. Я не зрозумів що з цим зробити. Відправити цей в сапорт? Щоб роздрукувати документ чи фото просто відправ його мені.",
	}
	SendToSupport = map[string]string{
		"uk": "Відправити в службу підтримки",
	}
	WasSentToSupport = map[string]string{
		"uk": "Відправлено вслужбу підтримки",
	}
)

func SelectLang(lang string) string {
	if lang == "uk" {
		return "uk"
	}
	return "uk"
}
