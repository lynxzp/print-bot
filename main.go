package main

import (
	"github.com/lynxzp/print-bot/cmd/printbot"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	printbot.Run()
}