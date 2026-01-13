package domain

import "context"

// ExpenseRepository defines the contract for storage operations.
// Any database (Postgres, MySQL, File, Memory) MUST implement these methods to work with our app.
// We use 'context' (ctx) for timeouts and cancellation - a strict industry requirement.
type ExpenseRepository interface {
	// Create saves a new expense and returns the ID or error
	Create(ctx context.Context, expense *Expense) error

	// GetAll retrieves all expenses
	GetAll(ctx context.Context) ([]Expense, error)

	// GetByID retrieves a specific expense
	GetByID(ctx context.Context, id int64) (*Expense, error)

	// Delete removes an expense by ID
	Delete(ctx context.Context, id int64) error
}

// ExpenseService defines the business logic contract.
// This is useful if we want to have multiple service implementations (rare but good for structure).
type ExpenseService interface {
	RegisterExpense(ctx context.Context, desc string, amount float64, category string) (*Expense, error)
	ListExpenses(ctx context.Context) ([]Expense, error)
	GetExpenseDetails(ctx context.Context, id int64) (*Expense, error)
	RemoveExpense(ctx context.Context, id int64) error
}
