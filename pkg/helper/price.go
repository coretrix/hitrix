package helper

import (
	"fmt"

	"github.com/bojanz/currency"
)

const (
	unit = 1000
)

type Price int64

func (c Price) Float() float64 {
	return float64(c) / unit
}

func (c Price) Units() int64 {
	return int64(c)
}

func (c Price) String() string {
	return fmt.Sprintf("%.2f", c.Float())
}

func (c Price) StringWithCurrency(currencySymbol string) string {
	return fmt.Sprintf("%.2f "+currencySymbol, c.Float())
}

func (c Price) StringByLocale(locale, inCurrency string) (string, error) {
	amount, err := currency.NewAmount(c.String(), inCurrency)

	if err != nil {
		return "", err
	}

	return currency.NewFormatter(currency.NewLocale(locale)).Format(amount), nil
}

func NewPrice(amount float64) Price {
	return Price(amount * unit)
}

func NewTotalPrice(amount float64, quantity uint64) Price {
	return NewPrice(float64(NewPrice(amount).Units()*int64(quantity)) / unit)
}
