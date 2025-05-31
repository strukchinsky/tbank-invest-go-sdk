package investapi

import (
	"fmt"

	"github.com/govalues/decimal"
)

func (q *Quotation) Decimal() decimal.Decimal {
	d, err := decimal.NewFromInt64(q.Units, int64(q.Nano), 9)
	if err != nil {
		panic(fmt.Sprintf("NewFromInt64(%v, %v, %v) failed: %v", q.Units, int64(q.Nano), 9, err))
	}
	return d
}

func NewQuotation(d decimal.Decimal) *Quotation {
	whole, frac, _ := d.Int64(9)
	return &Quotation{Units: whole, Nano: int32(frac)}
}
