package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	invoicepkg "github.com/werastine/Personalized_gift_bot.git/invoices"
	tele "gopkg.in/telebot.v4"
)

// divide code into different files to make an architecture
// I have to make a desicion what gift u want to send
// Then make to whom you want to send gift,
// You can choose to yourself/different person id
// end make auto send gift after purchase
// Да щас я напишу много хендлеров, но потом я постараюсь все перекинуть
// Очень сложно просто сейчас это все дается

// Sequence of everything
// Принять команду о покупке подарка
// Выбрать кому покупается подарок "Себе\другому"
// Принять то что напишет пользователь, заставив бота ждать сообщение от пользователя
// Выбрать подарок
// Выбор описания подарка тоже через принятие сообщ пользователя
// Отправить подарок

// Спросить у пользователя его

// there is UserID and his gift
// with/without
var userIDGift = make(map[int64]UserGiftToSend)
var userState = make(map[int64]string)

type UserGiftToSend struct {
	RecipientID     int64
	GiftID          string
	GiftDescription string
	ReciepUserName  string
	Ready           bool
}

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

	Cancle := tele.ReplyMarkup{}
	YesNotMenu := tele.ReplyMarkup{}
	inlineMenu := tele.ReplyMarkup{}
	keyboardBuyGift := tele.ReplyMarkup{
		OneTimeKeyboard: true,
		ResizeKeyboard:  true,
	}
	btnBear := inlineMenu.Data("🧸 15⭐️", "buy_bear")
	btnHeart := inlineMenu.Data("💝 15⭐️", "buy_heart")
	btnPresent := inlineMenu.Data("🎁 25⭐️", "buy_present")
	btnFlower := inlineMenu.Data("🌹 25⭐️", "buy_flower")
	btnBouqet := inlineMenu.Data("💐 50⭐️", "buy_bouqet")
	btnCake := inlineMenu.Data("🎂 50⭐️", "buy_cake")
	btnRocket := inlineMenu.Data("🚀 50⭐️", "buy_rocket")
	btnCup := inlineMenu.Data("🏆 100⭐️", "buy_cup")
	btnRing := inlineMenu.Data("💍 100⭐️", "buy_ring")
	btnChampagne := inlineMenu.Data("🍾 50⭐️", "buy_champagne")
	btnDiamond := inlineMenu.Data("💎 100⭐️", "buy_diamond")
	btnCancale := inlineMenu.Data("❌ Отмена", "buy_cancle")
	btnCancale1 := inlineMenu.Data("❌ Без Описания", "buy_cancle")
	btnNo := inlineMenu.Data("❌", "is_not_ok")
	btnYes := inlineMenu.Data("✅", "is_ok")
	//-------------------------------Replybuttons---------------------------//
	btnMyUser := keyboardBuyGift.Text("Подарок себе")
	btnOtherUser := keyboardBuyGift.User("Подарок другу", &tele.ReplyRecipient{
		ID:              1,
		RequestName:     tele.Flag(true),
		RequestUsername: tele.Flag(true),
	})
	//----------------------------------------------------------------------//

	inlineMenu.Inline(
		inlineMenu.Row(btnBear, btnHeart),
		inlineMenu.Row(btnPresent, btnFlower),
		inlineMenu.Row(btnBouqet, btnCake),
		inlineMenu.Row(btnRocket, btnCup),
		inlineMenu.Row(btnRing, btnChampagne),
		inlineMenu.Row(btnDiamond),
		inlineMenu.Row(btnCancale),
	)
	YesNotMenu.Inline(
		YesNotMenu.Row(btnYes, btnNo),
	)

	keyboardBuyGift.Reply(
		keyboardBuyGift.Row(btnMyUser, btnOtherUser),
	)
	Cancle.Inline(
		Cancle.Row(btnCancale1),
	)

	//---------------------------------BUTTON HANDLERS----------------------------//
	b.Handle(&btnCancale1, func(c tele.Context) error {
		defer delete(userState, c.Sender().ID)
		session := userIDGift[c.Sender().ID]
		session.GiftDescription = ""
		userIDGift[c.Sender().ID] = session
		c.Send("Переходим к оплате!")
		_, err := b.Send(c.Recipient(), invoicepkg.InvoiceHandler(session.GiftID))
		if err != nil {
			return err
		}
		return nil

	})

	b.Handle(&btnMyUser, func(c tele.Context) error {
		err := c.Send("Отлично, переходим к отправке", &tele.ReplyMarkup{RemoveKeyboard: true})
		if err != nil {
			return c.Send("Something went wrong")
		}
		c.Delete()
		session := userIDGift[c.Sender().ID]
		session.Ready = true
		session.RecipientID = c.Sender().ID
		userIDGift[c.Sender().ID] = session
		return c.Send("Выбери подарок для себя!", &inlineMenu)
	})

	// И по скольку у нас уже есть OnUserShared, то это надо будет как-то подсойденить чтобы можно было отправлять подарки по айдишнику
	// Человека. И бот будет почти готов!
	b.Handle(tele.OnUserShared, func(c tele.Context) error {
		shared := c.Message().UserShared
		if shared == nil {
			return c.Send("Ошибка при поиске пользоваетля!")
		}
		user := shared.Users[0]

		userMessage := c.Message()
		if err := b.Delete(userMessage); err != nil {
			return err
		}

		msg, err := b.Send(c.Recipient(), "Получен user!", &tele.ReplyMarkup{RemoveKeyboard: true})
		if err != nil {
			return err
		}
		b.Delete(msg)

		session := userIDGift[c.Sender().ID]
		session.RecipientID = user.UserID
		session.ReciepUserName = user.Username
		userIDGift[c.Sender().ID] = session

		return c.Send(fmt.Sprintf(
			"Это правильный получатель?\nИмя: %s\nUsername: @%s\nID: %d",
			user.FirstName,
			user.Username,
			user.UserID,
		), &YesNotMenu)
	})

	b.Handle(&btnYes, func(c tele.Context) error {
		if err := c.Respond(); err != nil {
			return err
		}
		session := userIDGift[c.Sender().ID]
		session.Ready = true
		userIDGift[c.Sender().ID] = session
		return c.Edit(fmt.Sprintf("Выбери подарок для %s", session.ReciepUserName), &inlineMenu)
	})

	b.Handle(&btnNo, func(c tele.Context) error {
		if err := c.Respond(); err != nil {
			return err
		}
		session := userIDGift[c.Sender().ID]
		session.Ready = false
		userIDGift[c.Sender().ID] = session
		c.Edit("Хорошо, попробуем заново!", &tele.ReplyMarkup{})
		return c.Send("Выбери получателя завово", &keyboardBuyGift)
	})

	b.Handle(&btnBear, func(c tele.Context) error {
		if err := c.Respond(); err != nil {
			return err
		}
		id := "5170233102089322756"
		session := userIDGift[c.Sender().ID]
		session.GiftID = id
		userIDGift[c.Sender().ID] = session
		c.Edit("Вы выбрали 🧸 за 15⭐️\n Теперь напиши описание подарка", &Cancle)
		//now we have invoice here but firstly we have to got the description
		userState[c.Sender().ID] = "wait_description"
		return nil
	})

	b.Handle(&btnHeart, func(c tele.Context) error {
		if err := c.Respond(); err != nil {
			return err
		}

		id := "5170145012310081615"
		session := userIDGift[c.Sender().ID]
		session.GiftID = id
		userIDGift[c.Sender().ID] = session

		c.Edit("Вы выбрали 💝 за 15 stars\n Теперь напиши описание подарка", &Cancle)
		userState[c.Sender().ID] = "wait_description"
		return nil
	})

	b.Handle(&btnPresent, func(c tele.Context) error {
		if err := c.Respond(); err != nil {
			return err
		}
		return c.Edit("Вы выбрали 🎁 за 15 stars, <b>Теперь переходим к оплате<b>")
	})
	//-----------------------------------------------------------------------------//

	b.Handle(tele.OnCallback, func(c tele.Context) error {
		return c.Send("Something went wrongg")
	})

	b.Handle("/start", func(c tele.Context) error {
		return c.Send("Привет! Это бот для авторских подарков\nНапиши /gift чтобы отправить подарок", &tele.ReplyMarkup{RemoveKeyboard: true})

	})

	b.Handle("/gift", func(c tele.Context) error {
		return c.Send("Привет, кому бы ты хотел оптравить подарок? ", &keyboardBuyGift)
	})

	b.Handle(tele.OnCheckout, func(c tele.Context) error {
		log.Println("Пришел запрос на подтверждение платежа")
		return c.Accept()

	})

	b.Handle(tele.OnPayment, func(c tele.Context) error {

		session, ok := userIDGift[c.Sender().ID]
		if !ok {
			return c.Send("Session is not found")
		}
		if session.RecipientID == 0 {
			return c.Send("Получатель не выбран")
		}
		if session.GiftID == "" {
			c.Send("Подарок не выбран")
		}
		recipient := &tele.User{ID: session.RecipientID}

		if err := b.SendGift(recipient, session.GiftID, session.GiftDescription); err != nil {
			log.Printf("SendGidt error: byer=%d, recipient=%d, gift=%s err=%v",
				c.Sender().ID, session.RecipientID, session.GiftID, err)
			c.Send("Оплата прошла но подарок не отправился")
		}

		log.Printf("Transaction info byer=%d, recipient=%d, gift=%s err=%v",
			c.Sender().ID, session.RecipientID, session.GiftID, err)
		return c.Send("Платеж обработан, отправляем подарок!")

	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		c.Edit("Принято!", &tele.ReplyMarkup{})
		msg := c.Message()
		if userState[c.Sender().ID] == "wait_description" {
			defer delete(userState, c.Sender().ID)
			session := userIDGift[c.Sender().ID]
			session.GiftDescription = msg.Text // there i will set gift desc. as user's message
			userIDGift[c.Sender().ID] = session

			_, err := b.Send(c.Recipient(), invoicepkg.InvoiceHandler(session.GiftID))
			if err != nil {
				return err
			}
			return nil
		}
		return b.Delete(msg)
	})

	b.Start()
}

// надо придумать логику чтобы все по перекидывать!
// Теперь надо дописать проект и добавить везде мютекс
