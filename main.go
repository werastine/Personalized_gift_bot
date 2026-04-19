package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	handlerspkg "github.com/werastine/Personalized_gift_bot.git/handlers"
	tele "gopkg.in/telebot.v4"
)

func main() {

	err := godotenv.Load("system/.env")
	if err != nil {
		log.Fatal("error loading env file")
	}

	pref := tele.Settings{
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	markup := handlerspkg.NewMarkupSet()
	storage := handlerspkg.NewStorage()

	app := &handlerspkg.App{
		Bot:     b,
		Markup:  markup,
		Storage: storage,
	}
	//In future i can divide RegisterHandlers code into few diff. groups
	app.RegisterHandlers()

	b.Start()
}
