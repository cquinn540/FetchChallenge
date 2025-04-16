package types

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

func (i *Item) UnmarshalJSON(b []byte) error {
	var item ItemRequest
	if err := json.Unmarshal(b, &item); err != nil {
		return err
	}

	price, err := decimal.NewFromString(item.Price)
	if err != nil {
		return err
	}

	*i = Item{
		ShortDescription: item.ShortDescription,
		Price:            price,
	}

	return nil
}

func (r *Receipt) UnmarshalJSON(b []byte) error {
	var receipt ReceiptRequest
	if err := json.Unmarshal(b, &receipt); err != nil {
		return err
	}

	purchaseDate, err := time.Parse(time.DateOnly, receipt.PurchaseDate)
	if err != nil {
		return err
	}

	purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
	if err != nil {
		return err
	}

	total, err := decimal.NewFromString(receipt.Total)
	if err != nil {
		return err
	}

	*r = Receipt{
		Retailer:     receipt.Retailer,
		PurchaseDate: purchaseDate,
		PurchaseTime: purchaseTime,
		Items:        receipt.Items,
		Total:        total,
	}

	return nil
}
