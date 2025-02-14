package Admin

import (
	"TgTaskBot/Config"
	"TgTaskBot/Log"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Admin struct {
	UserID int64
}

var SelectedAdmin *Admin

func (a *Admin) RequestAdminKey(bot *tgbotapi.BotAPI, message *tgbotapi.Message, updates tgbotapi.UpdatesChannel) {
	Log.FileLogger.Printf("INFO: User %s is requesting admin access", message.From.UserName)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Введите ключ доступа администратора:")
	row4 := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Вернуться в главное меню"),
	)
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(row4)

	if _, err := bot.Send(msg); err != nil {
		Log.FileLogger.Printf("ERROR: Failed to send admin key request message: %v", err)
		return
	}

	// Ожидаем ввода ключа администратора
	for update := range updates {
		if update.Message != nil {
			if update.Message.Text == "Вернуться в главное меню" {
				Log.FileLogger.Printf("INFO: User %s returned to the main menu from admin key request", message.From.UserName)
				return
			}

			if update.Message.Text == Config.AdminKey {
				Log.FileLogger.Printf("INFO: User %s successfully logged in as admin", message.From.UserName)
				if _, err := bot.Send(tgbotapi.NewMessage(message.Chat.ID, "✅ Вы вошли как администратор.")); err != nil {
					Log.FileLogger.Printf("ERROR: Failed to send admin login success message: %v", err)
				}
				a.ShowAdminOptions(bot, update.Message) // Показываем меню администратора
				return
			} else {
				Log.FileLogger.Printf("WARNING: User %s entered incorrect admin key: %s", message.From.UserName, update.Message.Text)
				if _, err := bot.Send(tgbotapi.NewMessage(message.Chat.ID, "⛔ Неверный ключ. Попробуйте снова.")); err != nil {
					Log.FileLogger.Printf("ERROR: Failed to send invalid admin key message: %v", err)
				}
			}
		}
	}
}

func (a *Admin) ShowAdminOptions(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	Log.FileLogger.Printf("INFO: Showing admin options to User %s", message.From.UserName)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите действие:")

	row0 := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Просмотр фотографий"),
	)
	row1 := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Просмотр файловой системы"),
	)
	row2 := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("просмотр файла логирования"),
	)
	row3 := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Добавить объект"),
		tgbotapi.NewKeyboardButton("Добавить должность"),
	)
	row4 := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Вернуться в главное меню"),
	)
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(row0, row1, row2, row3, row4)

	if _, err := bot.Send(msg); err != nil {
		Log.FileLogger.Printf("ERROR: Failed to send admin options message: %v", err)
	}
}
