package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"expense-tracker-v2/src/domain"

	_ "github.com/lib/pq" // Import the Postgres driver anonymously
)

// PostgresRepository implements domain.ExpenseRepository
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository is a factory to create our repository
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

// Create inserts a new expense into the database
func (r *PostgresRepository) Create(ctx context.Context, e *domain.Expense) error {
	query := `
		INSERT INTO expenses (description, amount, category, date)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	// standard context timeout handling is done by the caller (Service layer),
	// but QueryRowContext ensures we respect it.
	err := r.db.QueryRowContext(ctx, query, e.Description, e.Amount, e.Category, e.Date).Scan(&e.ID)
	if err != nil {
		return fmt.Errorf("failed to insert expense: %w", err)
	}
	return nil
}

// GetAll fetches all expenses
func (r *PostgresRepository) GetAll(ctx context.Context) ([]domain.Expense, error) {
	query := `SELECT id, description, amount, category, date FROM expenses`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query expenses: %w", err)
	}
	defer rows.Close() // Critical: always close rows to free connections

	var expenses []domain.Expense
	for rows.Next() {
		var e domain.Expense
		if err := rows.Scan(&e.ID, &e.Description, &e.Amount, &e.Category, &e.Date); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		expenses = append(expenses, e)
	}

	// check for errors that occurred during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return expenses, nil
}

// GetByID fetches a single expense
func (r *PostgresRepository) GetByID(ctx context.Context, id int64) (*domain.Expense, error) {
	query := `SELECT id, description, amount, category, date FROM expenses WHERE id = $1`

	var e domain.Expense
	err := r.db.QueryRowContext(ctx, query, id).Scan(&e.ID, &e.Description, &e.Amount, &e.Category, &e.Date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("expense not found")
		}
		return nil, fmt.Errorf("failed to get expense: %w", err)
	}
	return &e, nil
}

// Delete removes an expense
func (r *PostgresRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM expenses WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete expense: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("expense not found")
	}

	return nil
}
