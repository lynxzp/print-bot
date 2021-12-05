package printbot

import (
	"github.com/lynxzp/print-bot/pkg/bot"
	"github.com/lynxzp/print-bot/pkg/osenv"
	"log"
	"os"
)

func Run() {
	if len(os.Args) != 2 {
		log.Fatalln("second command line argument should be filename of json config")
	}
	cfg := readConfig(os.Args[1])
	osenv.Cfg = cfg.OSenv
	bot.Run(cfg.Bot)
}
