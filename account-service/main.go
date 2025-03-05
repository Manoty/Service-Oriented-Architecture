package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// Account struct
type Account struct {
	ID      int     `json:"id"`
	UserID  int     `json:"user_id"`
	Balance float64 `json:"balance"`
}

var (
	accounts = map[int]Account{
		1: {ID: 1, UserID: 1, Balance: 1000.0},
		2: {ID: 2, UserID: 2, Balance: 500.0},
	}
	mutex sync.Mutex // Prevents race conditions
)

// Get account balance
func GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	userID := 1 // Assume user ID is extracted from query
	account, exists := accounts[userID]
	if !exists {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(account)
}

// Update balance after a transaction
func UpdateBalanceHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserID int     `json:"user_id"`
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	account, exists := accounts[data.UserID]
	if !exists {
		mutex.Unlock()
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	// Apply balance change
	account.Balance += data.Amount
	accounts[data.UserID] = account
	mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Balance updated!"})
}

// Start the server
func main() {
	fmt.Println("✅ Account Service running on port 8083")

	http.HandleFunc("/account", GetAccountHandler)
	http.HandleFunc("/account/update", UpdateBalanceHandler)

	log.Fatal(http.ListenAndServe(":8083", nil))
}
