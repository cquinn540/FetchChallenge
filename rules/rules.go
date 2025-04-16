package rules

import (
	"FetchChallenge/types"
	"math"
	"strings"
	"unicode"

	"github.com/shopspring/decimal"
)

func CountAlphaNumeric(r *types.Receipt) int64 {
	var score int64 = 0
	for _, char := range r.Retailer {
		if (char >= '0' && char <= '9') || unicode.IsLetter(char) {
			score += 1
		}
	}

	return score
}

func IsTotalRound(r *types.Receipt) int64 {
	if r.Total.IsInteger() {
		return 50
	}
	return 0
}

func IsMultipleOfTwentyFiveCents(r *types.Receipt) int64 {
	pointTwoFive, _ := decimal.NewFromString("0.25")
	if r.Total.Mod(pointTwoFive).IsZero() {
		return 25
	}
	return 0
}

func FivePointsForEachTwoItems(r *types.Receipt) int64 {
	return int64(math.Floor(float64(len(r.Items))/2) * 5)
}

func itemDescriptionAndPrice(i *types.Item) int64 {
	length := len([]rune(strings.TrimSpace(i.ShortDescription)))

	if length%3 == 0 {
		pointTwo, _ := decimal.NewFromString("0.2")
		return i.Price.Mul(pointTwo).Ceil().BigInt().Int64()
	}

	return 0
}

func AllItemsDescriptionAndPrice(r *types.Receipt) int64 {
	var score int64 = 0
	for _, item := range r.Items {
		score += itemDescriptionAndPrice(&item)
	}

	return score
}

func IsPurchaseDateOdd(r *types.Receipt) int64 {
	if r.PurchaseDate.Day()%2 != 0 {
		return 6
	}
	return 0
}

func IsPurchaseBetweenTwoAndFour(r *types.Receipt) int64 {
	hour := r.PurchaseTime.Hour()
	if hour >= 14 && hour <= 16 {
		return 10
	}
	return 0
}

var ReceiptRules = [7]func(*types.Receipt) int64{
	CountAlphaNumeric,
	IsTotalRound,
	IsMultipleOfTwentyFiveCents,
	FivePointsForEachTwoItems,
	AllItemsDescriptionAndPrice,
	IsPurchaseDateOdd,
	IsPurchaseBetweenTwoAndFour,
}
