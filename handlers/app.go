package handlers

import (
	tele "gopkg.in/telebot.v4"
)

type App struct {
	Bot     *tele.Bot
	Storage *GlobalStorage
	Markup  *MarkupSet
}

type UserGiftToSend struct {
	RecipientID       int64
	GiftID            string
	GiftDescription   string
	ReciepUserName    string
	Ready             bool
	MessageIdToDelete tele.Editable
}

type MarkupSet struct {
	KeyboardBuyGift *tele.ReplyMarkup
	InlineMenu      *tele.ReplyMarkup
	YesNotMenu      *tele.ReplyMarkup
	Cancle          *tele.ReplyMarkup

	BtnBear      *tele.Btn
	BtnHeart     *tele.Btn
	BtnFlower    *tele.Btn
	BtnPresent   *tele.Btn
	BtnBouqet    *tele.Btn
	BtnCake      *tele.Btn
	BtnRocket    *tele.Btn
	BtnCup       *tele.Btn
	BtnRing      *tele.Btn
	BtnChampagne *tele.Btn
	BtnDiamond   *tele.Btn
	BtnCancale   *tele.Btn
	BtnCancale1  *tele.Btn
	BtnNo        *tele.Btn
	BtnYes       *tele.Btn
	//-------------------------------Replybuttons---------------------------//
	BtnMyUser    *tele.Btn
	BtnOtherUser *tele.Btn
	//etc
}
