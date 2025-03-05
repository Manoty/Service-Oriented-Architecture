package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Notification struct {
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
}

// Send a notification
func sendNotification(w http.ResponseWriter, r *http.Request) {
	var notification Notification
	json.NewDecoder(r.Body).Decode(&notification)

	fmt.Printf("📩 Notification sent to User %d: %s\n", notification.UserID, notification.Message)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "sent"})
}

func main() {
	http.HandleFunc("/notifications/send", sendNotification)

	fmt.Println("✅ Notification Service running on port 8083")
	http.ListenAndServe(":8083", nil)
}
