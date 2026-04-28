package handlers

import (
	"fmt"
	"log"

	invoicepkg "github.com/werastine/Personalized_gift_bot.git/invoices"
	tele "gopkg.in/telebot.v4"
)

func (a *App) RegisterHandlers() {

	a.Bot.Handle(a.Markup.BtnBear, func(c tele.Context) error {
		cb := c.Callback()
		if cb == nil || cb.Message == nil {
			return c.Send("Не можем получить айди сообщения")
		}
		if err := c.Respond(); err != nil {
			return err
		}
		id := "5170233102089322756"
		a.Storage.getSession(c.Sender().ID)
		a.Storage.setID(id, c.Sender().ID)
		a.Storage.DeleteMessage(c.Sender().ID, cb.Message)
		a.Storage.setState(c.Sender().ID, "wait_description")
		err := c.Edit("Вы выбрали 🧸 за 15⭐️\n Теперь напиши описание подарка", a.Markup.Cancle)
		if err != nil {
			return err
		}
		return nil
	})

	a.Bot.Handle(a.Markup.BtnHeart, func(c tele.Context) error {

		cb := c.Callback()
		if cb == nil || cb.Message == nil {
			return c.Send("Не можем получить айди сообщения")
		}
		if err := c.Respond(); err != nil {
			return err
		}
		id := "5170145012310081615"
		session := a.Storage.userIDGift[c.Sender().ID]
		session.GiftID = id
		session.MessageIdToDelete = cb.Message
		a.Storage.userState[c.Sender().ID] = "wait_description"
		a.Storage.userIDGift[c.Sender().ID] = session
		c.Edit("Вы выбрали 💝 за 15 stars\n Теперь напиши описание подарка", a.Markup.Cancle)
		return nil
	})

	a.Bot.Handle("/start", func(c tele.Context) error {
		return c.Send("Привет! ", &tele.ReplyMarkup{RemoveKeyboard: true})
	}) //Это бот для авторских подарков\nНапиши /gift чтобы отправить подарок

	a.Bot.Handle("/gift", func(c tele.Context) error {
		return c.Send("Привет, кому бы ты хотел оптравить подарок? ", a.Markup.KeyboardBuyGift)
	})

	a.Bot.Handle(a.Markup.BtnMyUser, func(c tele.Context) error {
		err := c.Send("Отлично, переходим к отправке", &tele.ReplyMarkup{RemoveKeyboard: true})
		if err != nil {
			return c.Send("Something went wrong")
		}
		c.Delete()
		session := a.Storage.userIDGift[c.Sender().ID]
		session.Ready = true
		session.RecipientID = c.Sender().ID
		a.Storage.userIDGift[c.Sender().ID] = session
		return c.Send("Выбери подарок для себя!", a.Markup.InlineMenu)
	})

	a.Bot.Handle(a.Markup.BtnCancale1, func(c tele.Context) error {
		defer delete(a.Storage.userState, c.Sender().ID)
		session := a.Storage.userIDGift[c.Sender().ID]
		session.GiftDescription = ""
		a.Storage.userIDGift[c.Sender().ID] = session
		errEdit := c.Edit("Переходим к оплате!", &tele.ReplyMarkup{})
		if errEdit != nil {
			log.Println("Error in editing message in BtnCancale1")
		}
		_, err := a.Bot.Send(c.Recipient(), invoicepkg.InvoiceHandler(session.GiftID))
		if err != nil {
			return err
		}
		return nil
	})

	a.Bot.Handle(a.Markup.BtnYes, func(c tele.Context) error {
		if err := c.Respond(); err != nil {
			return err
		}
		session := a.Storage.userIDGift[c.Sender().ID]
		session.Ready = true
		a.Storage.userIDGift[c.Sender().ID] = session
		return c.Edit(fmt.Sprintf("Выбери подарок для %s", session.ReciepUserName), a.Markup.InlineMenu)
	})

	a.Bot.Handle(a.Markup.BtnNo, func(c tele.Context) error {
		if err := c.Respond(); err != nil {
			return err
		}
		session := a.Storage.userIDGift[c.Sender().ID]
		session.Ready = false
		a.Storage.userIDGift[c.Sender().ID] = session
		c.Edit("Хорошо, попробуем заново!", &tele.ReplyMarkup{})
		return c.Send("Выбери получателя завово", a.Markup.KeyboardBuyGift)
	})

	//----------------------------On_Handlers----------------------------//

	a.Bot.Handle(tele.OnUserShared, func(c tele.Context) error {
		shared := c.Message().UserShared
		if shared == nil {
			return c.Send("Ошибка при поиске пользоваетля!")
		}
		user := shared.Users[0]

		userMessage := c.Message()
		if err := a.Bot.Delete(userMessage); err != nil {
			return err
		}

		msg, err := a.Bot.Send(c.Recipient(), "Получен user!", &tele.ReplyMarkup{RemoveKeyboard: true})
		if err != nil {
			return err
		}
		a.Bot.Delete(msg)

		session := a.Storage.userIDGift[c.Sender().ID]
		session.RecipientID = user.UserID
		session.ReciepUserName = user.Username
		a.Storage.userIDGift[c.Sender().ID] = session

		return c.Send(fmt.Sprintf(
			"Это правильный получатель?\nИмя: %s\nUsername: @%s\nID: %d",
			user.FirstName,
			user.Username,
			user.UserID,
		), a.Markup.YesNotMenu)
	})

	a.Bot.Handle(tele.OnCallback, func(c tele.Context) error {
		return c.Send("Something went wrongg")
	})

	a.Bot.Handle(tele.OnCheckout, func(c tele.Context) error {
		log.Println("Пришел запрос на подтверждение платежа")
		return c.Accept()
	})

	a.Bot.Handle(tele.OnPayment, func(c tele.Context) error {

		session, ok := a.Storage.userIDGift[c.Sender().ID]
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

		if err := a.Bot.SendGift(recipient, session.GiftID, session.GiftDescription); err != nil {
			log.Printf("SendGidt error: byer=%d, recipient=%d, gift=%s err=%v",
				c.Sender().ID, session.RecipientID, session.GiftID, err)
			c.Send("Оплата прошла но подарок не отправился")
		}

		log.Printf("Transaction info byer=%d, recipient=%d, gift=%s",
			c.Sender().ID, session.RecipientID, session.GiftID)
		return c.Send("Платеж обработан, отправляем подарок!")
	})

	a.Bot.Handle(tele.OnText, func(c tele.Context) error {
		msg := c.Message()
		if a.Storage.userState[c.Sender().ID] == "wait_description" {
			defer delete(a.Storage.userState, c.Sender().ID)
			session := a.Storage.userIDGift[c.Sender().ID]
			session.GiftDescription = msg.Text // there i will set gift desc. as user's message
			a.Storage.userIDGift[c.Sender().ID] = session
			_, errEdit := a.Bot.Edit(session.MessageIdToDelete, "Принято!", &tele.ReplyMarkup{})
			if errEdit != nil {
				log.Println("Error in editing message in Ontext")
			}

			_, err := a.Bot.Send(c.Recipient(), invoicepkg.InvoiceHandler(session.GiftID))
			if err != nil {
				return err
			}
			return nil
		}
		return a.Bot.Delete(msg)
	})

}
