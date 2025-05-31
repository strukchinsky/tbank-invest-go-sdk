package investapi

import (
	"fmt"

	"github.com/govalues/money"
)

func (m *MoneyValue) Amount() money.Amount {
	amount, err := money.NewAmountFromInt64(m.Currency, m.Units, int64(m.Nano), 9)
	if err != nil {
		panic(fmt.Sprintf("NewAmountFromInt64(%v, %v, %v, %v) failed with error %v", m.Currency, m.Units, int64(m.Nano), 9, err))
	}
	return amount
}

func NewMoneyValue(a money.Amount) *MoneyValue {
	whole, frac, _ := a.Int64(9)
	return &MoneyValue{Units: whole, Nano: int32(frac), Currency: a.Curr().Code()}
}
