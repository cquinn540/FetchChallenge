package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type ItemRequest struct {
	ShortDescription string
	Price            string
}

type ReceiptRequest struct {
	Retailer     string
	PurchaseDate string
	PurchaseTime string
	Items        []Item
	Total        string
}

type Item struct {
	ShortDescription string
	Price            decimal.Decimal
}

type Receipt struct {
	Retailer     string
	PurchaseDate time.Time
	PurchaseTime time.Time
	Items        []Item
	Total        decimal.Decimal
}

type ReceiptCreated struct {
	Id string `json:"id"`
}

type Score struct {
	Points string `json:"points"`
}
