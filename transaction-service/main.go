package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Transaction represents a fund transfer
type Transaction struct {
	FromUserID int     `json:"from_user_id"`
	ToUserID   int     `json:"to_user_id"`
	Amount     float64 `json:"amount"`
}

// API endpoint for Account Service
const accountServiceURL = "http://localhost:8083/account"

// Get account balance from Account Service
func getAccountBalance(userID int) (float64, error) {
	resp, err := http.Get(fmt.Sprintf("%s?id=%d", accountServiceURL, userID))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to fetch account balance")
	}

	var account struct {
		Balance float64 `json:"balance"`
	}
	err = json.NewDecoder(resp.Body).Decode(&account)
	if err != nil {
		return 0, err
	}

	return account.Balance, nil
}

// PerformTransaction handles transactions
func PerformTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check sender balance
	balance, err := getAccountBalance(transaction.FromUserID)
	if err != nil || balance < transaction.Amount {
		http.Error(w, "Insufficient funds", http.StatusBadRequest)
		return
	}

	// Update sender & receiver balances (mocked update)
	fmt.Printf("✅ Transaction successful: %d → %d ($%.2f)\n", transaction.FromUserID, transaction.ToUserID, transaction.Amount)

	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction completed!"})
}

func main() {
	fmt.Println("✅ Transaction Service running on port 8082")
	http.HandleFunc("/transactions/create", PerformTransaction)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
