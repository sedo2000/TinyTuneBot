package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		fmt.Fprintf(w, "Bot Error: %v", err)
		return
	}

	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK - TinyTune Bot is Listening! 🚀")
		return
	}

	if update.Message != nil && update.Message.Text == "/start" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to TinyTune! 🚀\nYour personal visualizer is ready.")
		
		// الرابط المباشر الذي طلبته
		directLink := "http://t.me/TinyTuneBot/visuals"

		// تعديل الزر ليكون رابط (URL) بدلاً من WebApp
		// تلجرام سيتعرف عليه تلقائياً ويفتحه كـ Mini App
		button := map[string]interface{}{
			"text": "🎵 Open TinyTune Visualizer",
			"url":  directLink,
		}

		keyboard := map[string]interface{}{
			"inline_keyboard": [][]interface{}{
				{button},
			},
		}

		keyboardBytes, _ := json.Marshal(keyboard)
		msg.ReplyMarkup = json.RawMessage(keyboardBytes)

		bot.Send(msg)
	}

	w.WriteHeader(http.StatusOK)
}
