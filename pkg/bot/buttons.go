package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lynxzp/print-bot/pkg/bot/text"
	"log"
	"strings"
)

// format of state type:
// state consist of list of params
// each param separated by '|' char
// order of params:
// 1 all params are visible or not (values: 1 and 0)
// 2 pages range (user readable range of printing pages, without spaces)
// 3 pages per list (available values: 1, 2, 4, 8, 16)
type state string

const (
	paramsCount = 3
)

func makeDocumentKeyboard(s state, lang string) tgbotapi.InlineKeyboardMarkup {
	bPage := tgbotapi.NewInlineKeyboardButtonData(text.Pages[lang]+" "+s.decodePages(), "changePages")
	return tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{bPage})
}

func (s *state) decodePages() string {
	s.shureValid()
	params := strings.Split(string(*s), "|")
	return params[2]
}

func (s *state) encodePages(p string) (written string) {
	s.shureValid()
	// cut bad chars
	p = strings.TrimSpace(p)
	notPermitted := strings.Trim(p, ",-1234567890")
	for _, c := range notPermitted {
		p = strings.ReplaceAll(p, string(c), "")
	}
	p = strings.Trim(p, notPermitted)
	// validate each range
	ranges := strings.Split(p, ",")
	var newRanges []string
	for _, r := range ranges {
		for (len(r) > 0) && (r[0] == '-') {
			r = r[1:len(r)]
		}
		for (len(r) > 0) && (r[len(r)-1] == '-') {
			r = r[0 : len(r)-1]
		}
		i1 := strings.Index(r, "-")
		i2 := strings.LastIndex(r, "-")

		if (len(r) > 0) && (i1 == i2) {
			newRanges = append(newRanges, r)
		}
	}
	p = strings.Join(newRanges, ",")
	// write
	params := strings.Split(string(*s), "|")
	params[2] = p
	*s = state(strings.Join(params, "|"))
	return p
}

func (s *state) shureValid() {
	if len(string(*s)) == 0 {
		s.new()
		return
	}
	params := strings.Split(string(*s), "|")
	if len(params) != paramsCount {
		s.logBroken()
		s.new()
		return
	}
	if (params[0] != "0") && (params[0] != "1") {
		s.logBroken()
		s.new()
	}
}

func (s *state) logBroken() {
	log.Println("EE broken state format:", s)
}

func (s *state) new() {
	*s = state("0||")
}
