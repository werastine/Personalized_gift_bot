package handlers

import tele "gopkg.in/telebot.v4"

func NewMarkupSet() *MarkupSet {

	cancle := tele.ReplyMarkup{}
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

	//-------------------------------ButtonForms---------------------------//
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
	cancle.Inline(
		cancle.Row(btnCancale1),
	)

	return &MarkupSet{
		//Forms
		InlineMenu:      &inlineMenu,
		YesNotMenu:      &YesNotMenu,
		KeyboardBuyGift: &keyboardBuyGift,
		Cancle:          &cancle,

		//buttons
		BtnBear:      &btnBear,
		BtnHeart:     &btnHeart,
		BtnPresent:   &btnPresent,
		BtnFlower:    &btnFlower,
		BtnBouqet:    &btnBouqet,
		BtnCake:      &btnCake,
		BtnRocket:    &btnRocket,
		BtnCup:       &btnCup,
		BtnRing:      &btnRing,
		BtnChampagne: &btnChampagne,
		BtnDiamond:   &btnDiamond,
		BtnCancale:   &btnCancale,
		BtnCancale1:  &btnCancale1,
		BtnNo:        &btnNo,
		BtnYes:       &btnYes,
		BtnMyUser:    &btnMyUser,
		BtnOtherUser: &btnOtherUser,
	}
}
