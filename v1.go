package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

// 1. DOMAIN: The shape of our data
type Expense struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Date        time.Time `json:"date"`
}

// Global variable to simulate our "database" in memory
var expenses []Expense

const fileName = "expenses.json"

// 2. STORAGE: Simple File I/O

// Load reads the JSON file into our 'expenses' slice
func loadExpenses() error {
	file, err := os.ReadFile(fileName)
	if os.IsNotExist(err) {
		return nil // File doesn't exist yet, that's fine (first run)
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(file, &expenses)
}

// Save writes our 'expenses' slice back to the JSON file
func saveExpenses() error {
	data, err := json.MarshalIndent(expenses, "", "  ") // Indent makes it readable
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, data, 0644)
}

// 3. LOGIC: Business operations
func addExpense(desc string, amount float64, category string) {
	// Auto-increment ID: typically last ID + 1
	id := 1
	if len(expenses) > 0 {
		id = expenses[len(expenses)-1].ID + 1
	}

	newExpense := Expense{
		ID:          id,
		Description: desc,
		Amount:      amount,
		Category:    category,
		Date:        time.Now(),
	}

	expenses = append(expenses, newExpense)

	if err := saveExpenses(); err != nil {
		fmt.Printf("Error saving data: %v\n", err)
	} else {
		fmt.Printf("Expense added successfully! ID: %d\n", id)
	}
}

func listExpenses() {
	fmt.Printf("%-5s | %-20s | %-10s | %-15s\n", "ID", "Description", "Amount", "Category")
	fmt.Println("------------------------------------------------------------")
	for _, e := range expenses {
		fmt.Printf("%-5d | %-20s | %-10.2f | %-15s\n", e.ID, e.Description, e.Amount, e.Category)
	}
}

func deleteExpense(id int) {
	index := -1
	// Find the index of the expense with the matching ID
	for i, e := range expenses {
		if e.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		fmt.Printf("Expense with ID %d not found\n", id)
		return
	}

	// The Standard Go "Delete from Slice" Trick:
	// Take everything BEFORE the index, and append everything AFTER the index.
	expenses = append(expenses[:index], expenses[index+1:]...)

	if err := saveExpenses(); err != nil {
		fmt.Printf("Error saving data: %v\n", err)
	} else {
		fmt.Printf("Expense deleted successfully! ID: %d\n", id)
	}
}

// 4. MAIN: entry point
func main() {
	// Define flags (commands)
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addDesc := addCmd.String("desc", "", "Description of the expense")
	addAmt := addCmd.Float64("amount", 0, "Amount of the expense")
	addCat := addCmd.String("cat", "General", "Category of the expense")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteID := deleteCmd.Int("id", 0, "ID of the expense to delete")

	// Load data at startup
	if err := loadExpenses(); err != nil {
		fmt.Printf("Error loading data: %v\n", err)
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("expected 'add' or 'list' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		if *addDesc == "" || *addAmt == 0 {
			fmt.Println("Please provide description and amount. Example: add -desc 'Lunch' -amount 50")
			return
		}
		addExpense(*addDesc, *addAmt, *addCat)

	case "delete":
		deleteCmd.Parse(os.Args[2:])
		if *deleteID == 0 {
			fmt.Println("Usage: delete -id 1")
			return
		}
		deleteExpense(*deleteID)

	case "list":
		listCmd.Parse(os.Args[2:])
		listExpenses()
	default:
		fmt.Println("expected 'add' or 'list' subcommands")
		os.Exit(1)
	}
}
