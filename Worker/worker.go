package Worker

import (
	"TgTaskBot/CategoryValidate"
	"TgTaskBot/Config"
	"TgTaskBot/ImgHandler"
	"TgTaskBot/Log"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Worker struct {
	Position   string
	ObjectName string
}

// StartPhotoUploadProcess запускает процесс загрузки фото
func (w *Worker) StartPhotoUploadProcess(bot *tgbotapi.BotAPI, message *tgbotapi.Message, updates tgbotapi.UpdatesChannel) {
	Log.FileLogger.Printf("INFO: Worker %s started the photo upload process", message.From.UserName)
	w.ShowObjectNameSelection(bot, message, updates)
}

// ShowObjectNameSelection отображает выбор объекта
func (w *Worker) ShowObjectNameSelection(bot *tgbotapi.BotAPI, message *tgbotapi.Message, updates tgbotapi.UpdatesChannel) {
	Log.FileLogger.Printf("DEBUG: Showing object selection for user: %s", message.From.UserName)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите объект:")

	// Создаем кнопки для объектов
	objectButtons := make([]tgbotapi.KeyboardButton, len(Config.Objects))
	for i, object := range Config.Objects {
		objectButtons[i] = tgbotapi.NewKeyboardButton(object)
	}

	// Группируем кнопки в строки по две
	var keyboardRows [][]tgbotapi.KeyboardButton
	row := []tgbotapi.KeyboardButton{}
	for i, button := range objectButtons {
		row = append(row, button)
		if (i+1)%2 == 0 {
			keyboardRows = append(keyboardRows, row)
			row = []tgbotapi.KeyboardButton{}
		}
	}

	// Добавляем оставшиеся кнопки, если они есть
	if len(row) > 0 {
		keyboardRows = append(keyboardRows, row)
	}

	// Кнопка для возврата в главное меню
	keyboardRows = append(keyboardRows, []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton("Вернуться в главное меню"),
	})

	// Применяем клавиатуру к сообщению
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(keyboardRows...)
	if _, err := bot.Send(msg); err != nil {
		Log.FileLogger.Printf("ERROR: Failed to send object selection message: %v", err)
		return
	}

	// Обработка выбора объекта через обновления
	for update := range updates {
		if update.Message != nil && update.Message.Chat.ID == message.Chat.ID {
			selectedObject := update.Message.Text
			if selectedObject == "Вернуться в главное меню" {
				Log.FileLogger.Printf("INFO: User %s returned to the main menu", message.From.UserName)
				ShowRoleSelection(bot, update.Message)
				return
			}
			if !CategoryValidate.IsObjectSelection(selectedObject) {
				Log.FileLogger.Printf("WARNING: User %s selected an invalid object: %s", message.From.UserName, selectedObject)
				if _, err := bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Выберите объект из списка")); err != nil {
					Log.FileLogger.Printf("ERROR: Failed to send invalid object message: %v", err)
				}
				continue
			}
			// Сохраняем выбранный объект
			w.ObjectName = selectedObject
			Log.FileLogger.Printf("INFO: User %s selected object: %s", message.From.UserName, selectedObject)

			// Шаг 2: Запрос позиции
			w.ShowPositionSelection(bot, update.Message, updates)
			break
		}
	}
}

// ShowPositionSelection отображает выбор должности
func (w *Worker) ShowPositionSelection(bot *tgbotapi.BotAPI, message *tgbotapi.Message, updates tgbotapi.UpdatesChannel) {
	Log.FileLogger.Printf("DEBUG: Showing position selection for user: %s", message.From.UserName)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите вашу должность:")

	// Создаем кнопки для позиций
	positionButtons := make([]tgbotapi.KeyboardButton, len(Config.Positions))
	for i, position := range Config.Positions {
		positionButtons[i] = tgbotapi.NewKeyboardButton(position)
	}

	// Группируем кнопки в строки по две
	var keyboardRows [][]tgbotapi.KeyboardButton
	row := []tgbotapi.KeyboardButton{}
	for i, button := range positionButtons {
		row = append(row, button)
		if (i+1)%2 == 0 {
			keyboardRows = append(keyboardRows, row)
			row = []tgbotapi.KeyboardButton{}
		}
	}

	// Добавляем оставшиеся кнопки, если они есть
	if len(row) > 0 {
		keyboardRows = append(keyboardRows, row)
	}

	// Кнопка для возврата в главное меню
	keyboardRows = append(keyboardRows, []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton("Вернуться в главное меню"),
	})

	// Применяем клавиатуру к сообщению
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(keyboardRows...)
	if _, err := bot.Send(msg); err != nil {
		Log.FileLogger.Printf("ERROR: Failed to send position selection message: %v", err)
		return
	}

	// Обработка выбора должности через обновления
	for update := range updates {
		if update.Message != nil && update.Message.Chat.ID == message.Chat.ID {
			selectedPosition := update.Message.Text
			if selectedPosition == "Вернуться в главное меню" {
				Log.FileLogger.Printf("INFO: User ID %d returned to the main menu", message.From.ID)
				ShowRoleSelection(bot, update.Message)
				return
			}
			if !CategoryValidate.IsPositionSelection(selectedPosition) {
				Log.FileLogger.Printf("WARNING: User ID %d selected an invalid position: %s", message.From.ID, selectedPosition)
				if _, err := bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Выберите должность из списка")); err != nil {
					Log.FileLogger.Printf("ERROR: Failed to send invalid position message: %v", err)
				}
				continue
			}
			// Сохраняем выбранную должность
			w.Position = selectedPosition
			Log.FileLogger.Printf("INFO: User ID %d selected position: %s", message.From.ID, selectedPosition)

			// Шаг 3: Запрос фото
			w.RequestPhoto(bot, update.Message)
			break
		}
	}
}

// RequestPhoto запрашивает фото у пользователя
func (w *Worker) RequestPhoto(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	Log.FileLogger.Printf("DEBUG: Requesting photo from user: %s", message.From.UserName)
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Вы выбрали %s для %s. \n Пожалуйста, отправьте фото.", w.Position, w.ObjectName))
	if _, err := bot.Send(msg); err != nil {
		Log.FileLogger.Printf("ERROR: Failed to send photo request message: %v", err)
	}
}

// HandlePhoto обрабатывает полученное фото
func (w *Worker) HandlePhoto(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	Log.FileLogger.Printf("INFO: Handling photo from user: %s", message.From.UserName)

	// Проверка, выбраны ли объект и должность
	if w.ObjectName == "" || w.Position == "" {
		Log.FileLogger.Printf("WARNING: User %s tried to upload a photo without selecting an object or position", message.From.UserName)
		if _, err := bot.Send(tgbotapi.NewMessage(message.Chat.ID, "Объект или роль не выбраны. Пожалуйста, начните выбор с начала")); err != nil {
			Log.FileLogger.Printf("ERROR: Failed to send error message: %v", err)
		}
		return
	}

	// Загрузка и сохранение фото
	filePath, err := ImgHandler.DownloadPhoto(bot, message, w.Position, w.ObjectName)
	if err != nil {
		Log.FileLogger.Printf("ERROR: Failed to download photo from user %d: %v", message.From.UserName, err)
		if _, err := bot.Send(tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("⛔ Ошибка: %v", err))); err != nil {
			Log.FileLogger.Printf("ERROR: Failed to send error message: %v", err)
		}
		return
	}

	// Логирование успешного сохранения фото
	Log.FileLogger.Printf("INFO: Photo from user %s saved; Path: %s", message.From.UserName, filePath)
	if _, err := bot.Send(tgbotapi.NewMessage(message.Chat.ID, "✅ Фотография успешно сохранена!")); err != nil {
		Log.FileLogger.Printf("ERROR: Failed to send success message: %v", err)
	}
}

// ShowRoleSelection отображает выбор роли
func ShowRoleSelection(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	Log.FileLogger.Printf("DEBUG: Showing role selection for user: %s", message.From.UserName)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите роль:")
	row1 := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Работник"),
		tgbotapi.NewKeyboardButton("Администратор"),
	)
	row2 := tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Вернуться в главное меню"))

	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(row1, row2)
	if _, err := bot.Send(msg); err != nil {
		Log.FileLogger.Printf("ERROR: Failed to send role selection message: %v", err)
	}
}
