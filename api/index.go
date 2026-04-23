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
		w.WriteHeader(http.StatusOK)
		return
	}

	var update tgbotapi.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "TinyTune Bot: Active 🚀")
		return
	}

	// معرف القناة (يجب أن يكون البوت مشرفاً فيها)
	const channelUsername = "@boxtoolls"
	const directLink = "http://t.me/TinyTuneBot/visuals"

	// 1. معالجة الضغط على زر التحقق (Callback Query)
	if update.CallbackQuery != nil {
		chatID := update.CallbackQuery.Message.Chat.ID
		userID := update.CallbackQuery.From.ID

		if update.CallbackQuery.Data == "verify_sub" {
			// التحقق من اشتراك المستخدم
			member, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
				ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
					SuperGroupUsername: channelUsername,
					UserID:             userID,
				},
			})

			if err == nil && (member.Status == "member" || member.Status == "administrator" || member.Status == "creator") {
				// إذا اشترك: احذف رسالة الاشتراك وأرسل رسالة الدخول
				bot.Send(tgbotapi.NewDeleteMessage(chatID, update.CallbackQuery.Message.MessageID))
				
				msg := tgbotapi.NewMessage(chatID, "Lost in the frequency.")
				button := map[string]interface{}{
					"text": "✨ DRIFT AWAY",
					"url":  directLink,
				}
				keyboard := map[string]interface{}{"inline_keyboard": [][]interface{}{{button}}}
				keyboardBytes, _ := json.Marshal(keyboard)
				msg.ReplyMarkup = json.RawMessage(keyboardBytes)
				bot.Send(msg)
			} else {
				// إذا لم يشترك: أرسل تنبيه
				callbackConfig := tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, "❌ Please subscribe to the channel first!")
				bot.Request(callbackConfig)
			}
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	// 2. معالجة أمر /start
	if update.Message != nil && update.Message.Text == "/start" {
		chatID := update.Message.Chat.ID
		userID := update.Message.From.ID

		// التحقق من الاشتراك قبل إظهار أي شيء
		member, err := bot.GetChatMember(tgbotapi.GetChatMemberConfig{
			ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
				SuperGroupUsername: channelUsername,
				UserID:             userID,
			},
		})

		if err == nil && (member.Status == "member" || member.Status == "administrator" || member.Status == "creator") {
			// مشترك بالفعل: أرسل الرسالة الأصلية مباشرة
			msg := tgbotapi.NewMessage(chatID, "Lost in the frequency.")
			button := map[string]interface{}{"text": "✨ DRIFT AWAY", "url": directLink}
			keyboard := map[string]interface{}{"inline_keyboard": [][]interface{}{{button}}}
			keyboardBytes, _ := json.Marshal(keyboard)
			msg.ReplyMarkup = json.RawMessage(keyboardBytes)
			bot.Send(msg)
		} else {
			// غير مشترك: أرسل رسالة الاشتراك الإجباري
			msg := tgbotapi.NewMessage(chatID, "⚠️ Access Denied!\n\nPlease subscribe to our channel to use the bot.")
			
			btnSub := map[string]interface{}{"text": "📢 Subscribe", "url": "https://t.me/boxtoolls"}
			btnVerify := map[string]interface{}{"text": "✅ Verify Subscription", "callback_data": "verify_sub"}
			
			keyboard := map[string]interface{}{
				"inline_keyboard": [][]interface{}{{btnSub}, {btnVerify}},
			}
			
			keyboardBytes, _ := json.Marshal(keyboard)
			msg.ReplyMarkup = json.RawMessage(keyboardBytes)
			bot.Send(msg)
		}
	}

	w.WriteHeader(http.StatusOK)
}
