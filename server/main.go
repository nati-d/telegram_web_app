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
		Token:  "YOUR_BOT_TOKEN_HERE",
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
			URL: "https://your-vercel-url.vercel.app/",
		})
		markup.Reply(markup.Row(btn))
		return c.Send("Tap below to open the WebApp 👇", markup)
	})

	// Run the bot in a separate goroutine
	go bot.Start()

	// REST API endpoint to receive data from React
	http.HandleFunc("/api/submit", func(w http.ResponseWriter, r *http.Request) {
		var data MessageData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Received from user %s (%d): %s", data.Username, data.UserID, data.Message)

		// Reply to user via Telegram
		_, err := bot.Send(&telebot.User{ID: data.UserID}, "✅ Received your message: "+data.Message)
		if err != nil {
			log.Println("Error sending message:", err)
		}

		w.WriteHeader(http.StatusOK)
	})

	log.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
