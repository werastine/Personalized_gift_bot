package invoices

import (
	tele "gopkg.in/telebot.v4"
)

func InvoiceHandler(id string) *tele.Invoice {
	switch id {
	case "5170233102089322756":
		return SendBearInvoice()
	case "5170145012310081615":
		return SendHeartInvoice()
	}

	return nil
}

// -----------------------------------------------------------------------------//
func SendBearInvoice() *tele.Invoice {

	price := tele.Price{
		Label:  "Bear",
		Amount: 15,
	}

	InvoiceBear := tele.Invoice{
		Title: "🧸 Подарок",
		// тут когда я сделаю архитектуру будет писать кому отправляется и описание подарка
		Description: "Порадует любого!",
		Currency:    "XTR",
		Payload:     "Bear_001",
		Prices:      []tele.Price{price},
	}
	return &InvoiceBear
}

//-----------------------------------------------------------------------------//

func SendHeartInvoice() *tele.Invoice {
	price := tele.Price{
		Label:  "Bear",
		Amount: 15,
	}

	InvoiceHeart := tele.Invoice{
		Title:       "💝 Подарок",
		Description: "Порадует любого!",
		Currency:    "XTR",
		Payload:     "Bear_001",
		Prices:      []tele.Price{price},
	}
	return &InvoiceHeart
}

//-----------------------------------------------------------------------------//

func SendPresentInvoice() *tele.Invoice {
	price := tele.Price{
		Label:  "Bear",
		Amount: 25,
	}

	InvoiceHeart := tele.Invoice{
		Title:       "🎁 Подарок",
		Description: "Порадует любого!",
		Currency:    "XTR",
		Payload:     "Bear_001",
		Prices:      []tele.Price{price},
	}
	return &InvoiceHeart
}
