package domain

import (
	"time"
)

// Expense represents a single financial record in our system.
// We use standard JSON tags for API responses and DB tags for potential database mapping.
type Expense struct {
	ID          int64     `json:"id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Date        time.Time `json:"date"`
}

// NewExpense is a factory function to create a validated expense (Optional but good practice)
func NewExpense(desc string, amount float64, category string) *Expense {
	return &Expense{
		Description: desc,
		Amount:      amount,
		Category:    category,
		Date:        time.Now(),
	}
}
