package investapi

import "github.com/govalues/money"

func (m *MoneyValue) Amount() (money.Amount, error) {
	return money.NewAmountFromInt64(m.Currency, m.Units, int64(m.Nano), 9)
}

func NewMoneyValue(a money.Amount) *MoneyValue {
	whole, frac, _ := a.Int64(9)
	return &MoneyValue{Units: whole, Nano: int32(frac), Currency: a.Curr().Code()}
}
