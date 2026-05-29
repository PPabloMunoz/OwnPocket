package model

// Amount represents monetary value in cents (int64)
// Example: 1250 = $12.50
type Amount int64

func (a Amount) Cents() int64 {
	return int64(a)
}

func (a Amount) ToFloat() float64 {
	return float64(a) / 100
}

func NewAmountFromCents(cents int64) Amount {
	return Amount(cents)
}

func NewAmountFromFloat(dollars float64) Amount {
	return Amount(int64(dollars * 100))
}
