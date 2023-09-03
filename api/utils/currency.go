package utils

// IntToFloatCurrency converts an integer currency value (in cents) to a floating-point value (in dollars).
func IntToFloatCurrency(amountInCents int64) float64 {
	return float64(amountInCents) / 100.0
}

// FloatToIntCurrency converts a floating-point currency value (in dollars) to an integer value (in cents).
func FloatToIntCurrency(amountInDollars float64) int64 {
	return int64(amountInDollars * 100.0)
}
