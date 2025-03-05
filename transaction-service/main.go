package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Transaction struct {
	FromUserID int     `json:"from_user_id"`
	ToUserID   int     `json:"to_user_id"`
	Amount     float64 `json:"amount"`
}

var transactions []Transaction

// Notify function to call notification-service
func sendNotification(userID int, message string) {
	notification := map[string]interface{}{
		"user_id":  userID,
		"message":  message,
	}

	notificationJSON, _ := json.Marshal(notification)
	_, err := http.Post("http://localhost:8083/notifications/send", "application/json", bytes.NewBuffer(notificationJSON))
	if err != nil {
		fmt.Println("❌ Failed to send notification:", err)
	}
}

// Process a new transaction and send a notification
func createTransaction(w http.ResponseWriter, r *http.Request) {
	var newTransaction Transaction
	json.NewDecoder(r.Body).Decode(&newTransaction)

	transactions = append(transactions, newTransaction)

	fmt.Printf("💰 Transaction: %d -> %d | Amount: %.2f\n", newTransaction.FromUserID, newTransaction.ToUserID, newTransaction.Amount)

	// Send notification
	message := fmt.Sprintf("You received a payment of %.2f", newTransaction.Amount)
	sendNotification(newTransaction.ToUserID, message)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTransaction)
}

func main() {
	http.HandleFunc("/transactions/create", createTransaction)

	fmt.Println("✅ Transaction Service running on port 8082")
	http.ListenAndServe(":8082", nil)
}
