package investgo

import (
	"math"

	pb "github.com/floatdrop/tbank-invest-go-sdk"
)

// billion used for converting floats to Units and Nano parts
// Fractional part of float should be multiplied by billion to get Nanos and vice versa
const billion int64 = 1_000_000_000

// QuotationToFloat converts Quotation type to float64 with proper scaling
func QuotationToFloat(q *pb.Quotation) float64 {
	return float64(q.Units) + float64(q.Nano)/float64(billion)
}

// FloatToQuotation converts float64 to Quotation type with proper scaling
func FloatToQuotation(f float64) *pb.Quotation {
	openUnits, openNanos := math.Modf(f)
	return &pb.Quotation{Units: int64(openUnits), Nano: int32(openNanos * float64(billion))}
}
