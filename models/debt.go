package models

// Debt model
type Debt struct {
	ID       int
	Customer string
	Phone    string
	Amount   float64
	Balance  float64
}
