package service

import (
	"context"
	"errors"
	"time"

	"expense-tracker-v2/src/domain"
)

// expenseService implements domain.ExpenseService
// It holds a reference to the Repository Interface, not the concrete struct.
type expenseService struct {
	repo domain.ExpenseRepository
}

// NewExpenseService is the constructor
func NewExpenseService(repo domain.ExpenseRepository) domain.ExpenseService {
	return &expenseService{
		repo: repo,
	}
}

// RegisterExpense validates input and creates the expense
func (s *expenseService) RegisterExpense(ctx context.Context, desc string, amount float64, category string) (*domain.Expense, error) {
	// 1. Business Validation Logic
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}
	if desc == "" {
		return nil, errors.New("description cannot be empty")
	}

	// 2. Prepare the data
	newExpense := &domain.Expense{
		Description: desc,
		Amount:      amount,
		Category:    category,
		Date:        time.Now(),
	}

	// 3. Persist using the repository
	err := s.repo.Create(ctx, newExpense)
	if err != nil {
		return nil, err
	}

	return newExpense, nil
}

// ListExpenses passes the call to the repository
func (s *expenseService) ListExpenses(ctx context.Context) ([]domain.Expense, error) {
	return s.repo.GetAll(ctx)
}

// GetExpenseDetails by id
func (s *expenseService) GetExpenseDetails(ctx context.Context, id int64) (*domain.Expense, error) {
	return s.repo.GetByID(ctx, id)
}

// RemoveExpense passes the call to the repository
func (s *expenseService) RemoveExpense(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
