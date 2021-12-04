package printbot

import (
	"encoding/json"
	"github.com/lynxzp/print-bot/pkg/bot"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Bot bot.Config
}


func readConfig(filename string) (cfg Config){
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln("Can't read config")
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err!=nil {
		log.Fatalln("Can't read config")
	}
	err = json.Unmarshal(b, &cfg)
	if err!=nil {
		log.Fatalln(err)
	}
	return cfg
}