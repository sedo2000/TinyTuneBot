package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// جلب التوكن من Environment Variables في فيرسل
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		// إذا كان هناك خطأ في التوكن سيظهر هنا في الـ Logs
		fmt.Fprintf(w, "Bot Error: %v", err)
		return
	}

	var update tgbotapi.Update
	// استقبال البيانات من تلجرام
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		// عند فتح الرابط من المتصفح يدوياً سيصل هنا لأنه لا توجد بيانات JSON
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK - TinyTune Bot is Listening! 🚀")
		return
	}

	if update.Message != nil && update.Message.Text == "/start" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to TinyTune! 🚀\nYour personal visualizer is ready.")
		
		// بناء رابط الـ WebApp ديناميكياً
		webAppURL := "https://" + r.Host + "/"

		// استخدام Raw JSON لتجنب مشاكل إصدار المكتبة في فيرسل
		button := map[string]interface{}{
			"text": "🎵 Open TinyTune Visualizer",
			"web_app": map[string]string{
				"url": webAppURL,
			},
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
