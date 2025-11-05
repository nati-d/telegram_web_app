package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/telebot.v3"
)

type MessageData struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	Message   string `json:"message"`
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using system environment variables")
	}

	botToken := os.Getenv("BOT_TOKEN")
	webAppURL := os.Getenv("WEBAPP_URL")
	port := os.Getenv("PORT")

	if botToken == "" {
		log.Fatal("❌ BOT_TOKEN not set in .env")
	}

	pref := telebot.Settings{
		Token:  botToken,
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
			URL: webAppURL,
		})
		markup.Reply(markup.Row(btn))
		return c.Send("Tap below to open the WebApp 👇", markup)
	})

	// Run the bot in a separate goroutine
	go bot.Start()

	// REST API endpoint to receive data from React
	http.HandleFunc("/api/submit", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		var data MessageData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("📩 Received from user %s (%d): %s", data.Username, data.UserID, data.Message)

		_, err := bot.Send(&telebot.User{ID: data.UserID}, "✅ Received your message: "+data.Message)
		if err != nil {
			log.Println("❌ Error sending message:", err)
		}

		w.WriteHeader(http.StatusOK)
	})

	log.Printf("🚀 Server running on http://localhost:%s", port)
	http.ListenAndServe(":"+port, nil)
}
