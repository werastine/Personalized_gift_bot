package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
)

// divide code into different files to make an architecture
// I have to make a desicion what gift u want to send
// Then make to whom you want to send gift,
// You can choose to yourself/different person id
// end make auto send gift after purchase

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

	price := tele.Price{
		Label:  "Bear",
		Amount: 15,
	}

	InvoiceBear := tele.Invoice{
		Title:       "🧸 Подарок",
		Description: "Для selected user, with/without descrtiprion",
		Currency:    "XTR",
		Payload:     "Bear_001",
		Prices:      []tele.Price{price},
	}
	bear := tele.Gift{
		ID:        "5170233102089322756",
		StarCount: 15,
	}

	b.Handle("/start", func(c tele.Context) error {
		return c.Send("Hello write /gift")

	})

	b.Handle("/gift", func(c tele.Context) error {
		return b.SendGift(c.Recipient(), bear.ID)

	})

	b.Handle(tele.OnCheckout, func(c tele.Context) error {
		log.Println("Пришел запрос на подтверждение платежа")
		return c.Accept()

	})

	b.Handle(tele.OnPayment, func(c tele.Context) error {
		// there i will make switch case, if Bear - send bear, if diamond send diamond
		return c.Send("Платеж обработан, отправляем подарок!")

	})

	b.Handle("/pay", func(c tele.Context) error {
		b.Send(c.Recipient(), &InvoiceBear)
		return nil
	})
	b.Start()
}

// b.Handle(tele.OnText, func(c tele.Context) error {
// 	msg := c.Text()
// 	name := c.Chat().FirstName
// 	return c.Send(("Your name is " + name + "Your text was:" + msg))
// })
