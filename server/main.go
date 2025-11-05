package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"gopkg.in/telebot.v3"
)

type MessageData struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	Message   string `json:"message"`
}

func main() {
	pref := telebot.Settings{
		Token:  "8029064537:AAGFVMrZiuoKpEOD7TnPLJ7zlxDT3IonJlM",
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Handle /start to show the WebApp button
	bot.Handle("/start", func(c telebot.Context) error {
		markup := &telebot.ReplyMarkup{ResizeKeyboard: true}
		btn := markup.WebApp("Open WebApp", &telebot.WebApp{
			URL: "https://telegram-web-app-five-theta.vercel.app/",
		})
		markup.Reply(markup.Row(btn))
		return c.Send("Tap below to open the WebApp 👇", markup)
	})

	// Run the bot in a separate goroutine
	go bot.Start()

	// REST API endpoint to receive data from React
http.HandleFunc("/api/submit", func(w http.ResponseWriter, r *http.Request) {
    // Allow requests from your Vite dev server
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // Handle preflight OPTIONS request
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }

    var data MessageData
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    log.Printf("Received from user %s (%d): %s", data.Username, data.UserID, data.Message)

    // Send Telegram reply
    _, err := bot.Send(&telebot.User{ID: data.UserID}, "✅ Received your message: "+data.Message)
    if err != nil {
        log.Println("Error sending message:", err)
    }

    w.WriteHeader(http.StatusOK)
})

	log.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
