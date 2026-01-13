package transport

import (
	"encoding/json"
	"expense-tracker-v2/src/domain"
	"net/http"
	"strconv"
	"strings"
)

// Handler holds the service interface
type Handler struct {
	Service domain.ExpenseService
}

// NewHandler creates our HTTP handler
func NewHandler(service domain.ExpenseService) *Handler {
	return &Handler{
		Service: service,
	}
}

// CreateExpenseHandler handles POST /expenses
func (h *Handler) CreateExpenseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Define a temporary struct for incoming JSON payload
	type request struct {
		Description string  `json:"description"`
		Amount      float64 `json:"amount"`
		Category    string  `json:"category"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the service
	expense, err := h.Service.RegisterExpense(r.Context(), req.Description, req.Amount, req.Category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Send response
	// FIX: Set Header first, then Status Code, then Body
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Returns 201
	json.NewEncoder(w).Encode(expense)
}

// ListExpensesHandler handles GET /expenses
func (h *Handler) ListExpensesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	expenses, err := h.Service.ListExpenses(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch expenses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

func (h *Handler) GetExpenseByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path (Standard REST)
	idStr := strings.TrimPrefix(r.URL.Path, "/expenses/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Call Service
	// Note: We need to ensure your Service Interface has GetByID.
	// If we missed exporting it in the Service Interface earlier, use the Repo directly or add it to Service.
	// For now, assuming your Service interface has a 'GetDetail' or similar.
	// If not, let's add it (see step below).

	// *Assuming you might not have added GetByID to the Service Layer yet:*
	// Let's assume we add it now.
	expense, err := h.Service.GetExpenseDetails(r.Context(), id)
	if err != nil {
		http.Error(w, "Expense not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expense)
}

func (h *Handler) DeleteExpenseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// EXTRACT ID FROM URL PATH
	// Path is "/expenses/123"
	// We trim the prefix "/expenses/" to get "123"
	idStr := strings.TrimPrefix(r.URL.Path, "/expenses/")

	// Safety check: if they just called "/expenses/", idStr is empty
	if idStr == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Call service
	err = h.Service.RemoveExpense(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Expense deleted"))
}
