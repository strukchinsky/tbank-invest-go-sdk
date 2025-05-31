package investapi

import (
	"testing"

	"github.com/govalues/decimal"
	"google.golang.org/protobuf/proto"
)

func TestDecimal(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		in  *Quotation
		out decimal.Decimal
	}{
		"100.0 is 100.0": {in: &Quotation{Units: 100, Nano: 0}, out: decimal.MustParse("100.0")},
		"0.50":           {in: &Quotation{Units: 0, Nano: 500000000}, out: decimal.MustParse("0.5")},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			d := test.in.Decimal()
			if !d.Equal(test.out) {
				t.Errorf("%v should be equal to %v", d, test.out)
			}
		})
	}
}

func TestNewQuotation(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		in  decimal.Decimal
		out *Quotation
	}{
		"100$ is 100$": {in: decimal.MustParse("100.0"), out: &Quotation{Units: 100, Nano: 0}},
		"50 cents":     {in: decimal.MustParse("0.5"), out: &Quotation{Units: 0, Nano: 500000000}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			value := NewQuotation(test.in)

			if !proto.Equal(value, test.out) {
				t.Errorf("{%v} should be equal to {%v}", value, test.out)
			}
		})
	}
}
