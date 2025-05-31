package investapi

import (
	"testing"

	"github.com/govalues/money"
	"google.golang.org/protobuf/proto"
)

type MoneyAmountWithError struct {
	amount money.Amount
	err    error
}

func NewAmountWithError(curr string, value int64, scale int) MoneyAmountWithError {
	a, e := money.NewAmount(curr, value, scale)
	return MoneyAmountWithError{amount: a, err: e}
}

func TestAmount(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		in  *MoneyValue
		out MoneyAmountWithError
	}{
		"100$ is 100$": {in: &MoneyValue{Currency: "USD", Units: 100, Nano: 0}, out: NewAmountWithError("USD", 100, 0)},
		"50 cents":     {in: &MoneyValue{Currency: "USD", Units: 0, Nano: 500000000}, out: NewAmountWithError("USD", 50, 2)},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			amount, err := test.in.Amount()
			if err != test.out.err {
				t.Errorf(`Amount() should return error %v but got %v`, test.out.err, err)
			}

			equal, err := amount.Equal(test.out.amount)
			if !equal || err != nil {
				t.Errorf("%v should be equal to %v (error is %v)", amount, test.out.amount, err)
			}
		})
	}
}

func TestNewMoneyValue(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		in  money.Amount
		out *MoneyValue
	}{
		"100$ is 100$": {in: money.MustParseAmount("USD", "100.0"), out: &MoneyValue{Currency: "USD", Units: 100, Nano: 0}},
		"50 cents":     {in: money.MustParseAmount("USD", "0.5"), out: &MoneyValue{Currency: "USD", Units: 0, Nano: 500000000}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			value := NewMoneyValue(test.in)

			if !proto.Equal(value, test.out) {
				t.Errorf("{%v} should be equal to {%v}", value, test.out)
			}
		})
	}
}
