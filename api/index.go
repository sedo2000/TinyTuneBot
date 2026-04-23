package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// 1. جلب التوكن من إعدادات البيئة (أمان أكثر)
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		http.Error(w, "Bot initialization error", http.StatusInternalServerError)
		return
	}

	// 2. فك تشفير البيانات القادمة من تلجرام
	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Failed to decode update", http.StatusBadRequest)
		return
	}

	// 3. معالجة الرسائل
	if update.Message != nil {
		chatID := update.Message.Chat.ID
		text := update.Message.Text

		var msg tgbotapi.MessageConfig

		if text == "/start" {
			// إعداد زر الـ Mini App
			webApp := &tgbotapi.WebAppInfo{
				URL: "https://your-tinytune-frontend.vercel.app", // رابط موقعك
			}
			button := tgbotapi.InlineKeyboardButton{
				Text:   "🎵 Open TinyTune",
				WebApp: webApp,
			}
			keyboard := tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{button})

			msg = tgbotapi.NewMessage(chatID, "Welcome to TinyTune! 🚀\nYour personal visualizer is ready.")
			msg.ReplyMarkup = keyboard
		} else {
			msg = tgbotapi.NewMessage(chatID, "I only respond to /start for now!")
		}

		bot.Send(msg)
	}

	// 4. الرد على تلجرام بنجاح (200 OK) لكي لا يعيد إرسال الرسالة
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}
