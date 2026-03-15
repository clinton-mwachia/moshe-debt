package models

// Payment model
type Payment struct {
	ID        int
	Customer  string
	Amount    float64
	Balance   float64
	CreatedAt string
}
