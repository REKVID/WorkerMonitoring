package Admin

import (
	"TgTaskBot/Config"
	"TgTaskBot/pkg/handlers"
	Log "TgTaskBot/pkg/logger"
	"log"
	"os"
	"path/filepath"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (a *Admin) AddObject(bot *tgbotapi.BotAPI, message *tgbotapi.Message, updates tgbotapi.UpdatesChannel) {
	Log.FileLogger.Printf("INFO: Admin %s requested to add an object", message.From.UserName)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Введите название нового объекта:")
	bot.Send(msg)

	for update := range updates {
		if update.Message != nil && update.Message.Text != "" {
			newObject := update.Message.Text
			Config.Objects = append(Config.Objects, newObject)
			err := Config.SaveConfig("./Config/configuration.json")
			if err != nil {
				Log.FileLogger.Printf("ERROR: Couldn't update config file: %v", err)
				bot.Send(tgbotapi.NewMessage(message.Chat.ID, "⛔ Не удалось обновить файл конфигурации."))
				return
			}
			Log.FileLogger.Printf("INfO: Admin %s added new object: %s", message.From.UserName, newObject)
			bot.Send(tgbotapi.NewMessage(message.Chat.ID, "✅ Объект успешно добавлен!"))
			break
		}
	}
}

func (a *Admin) AddPosition(bot *tgbotapi.BotAPI, message *tgbotapi.Message, updates tgbotapi.UpdatesChannel) {
	Log.FileLogger.Printf("INFO: Admin %d requested to add a position", message.From.ID)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Введите название новой позиции:")
	bot.Send(msg)

	for update := range updates {
		if update.Message != nil && update.Message.Text != "" {
			newPosition := update.Message.Text
			Config.Positions = append(Config.Positions, newPosition)
			err := Config.SaveConfig("./Config/configuration.json")
			if err != nil {
				Log.FileLogger.Printf("ERROR: Couldn't update config file: %v", err)
				bot.Send(tgbotapi.NewMessage(message.Chat.ID, "⛔ Не удалось обновить файл конфигурации."))
				return
			}
			Log.FileLogger.Printf("INFO: Admin %d added new position: %s", message.From.ID, newPosition)
			bot.Send(tgbotapi.NewMessage(message.Chat.ID, "✅ Позиция успешно добавлена!"))
			break
		}
	}
}

func (a *Admin) GetFilteredPhotos(bot *tgbotapi.BotAPI, message *tgbotapi.Message, updates tgbotapi.UpdatesChannel) {

	selectedPosition := a.selectPosition(bot, message, updates)
	if selectedPosition == "" {
		return
	}

	selectedObject := a.selectObject(bot, message, updates)
	if selectedObject == "" {
		return
	}

	selectedDate := a.selectDate(bot, message, updates)
	if selectedDate == "" {
		return
	}

	a.sendPhotos(bot, message, selectedPosition, selectedObject, selectedDate)
}

func (a *Admin) ListFileSystem(bot *tgbotapi.BotAPI, message *tgbotapi.Message, baseDir string) {

	// Проверяем, существует ли BaseDir
	if _, err := os.Stat(baseDir); os.IsNotExist(err) {
		Log.FileLogger.Printf("ERROR: Directory %s does not exist: %v", baseDir, err)
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "⛔ Директория "+baseDir+" не существует."))
		return
	}

	// Проверяем, является ли BaseDir директорией
	if info, err := os.Stat(baseDir); err != nil {
		Log.FileLogger.Printf("ERROR: Failed to access directory %s: %v", baseDir, err)
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "⛔ Ошибка при доступе к директории."))
		return
	} else if !info.IsDir() {
		Log.FileLogger.Printf("ERROR: %s is not a directory", baseDir)
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "⛔ "+baseDir+" не является директорией."))
		return
	}

	// Строим дерево файловой системы
	tree := buildFileTree(baseDir, "", 0)
	if tree == "" {
		Log.FileLogger.Printf("WARNING: Directory %s is empty or inaccessible", baseDir)
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "⚠ Директория пуста или недоступна."))
		return
	}

	// Отправляем сообщение с MarkdownV2
	msg := tgbotapi.NewMessage(message.Chat.ID, "Файловая система:\n```\n"+tree+"\n```")
	msg.ParseMode = "MarkdownV2" // Включаем MarkdownV2
	_, err := bot.Send(msg)
	if err != nil {
		Log.FileLogger.Printf("ERROR: Failed to send file system tree message: %v", err)
	}
}

func (a *Admin) SendLogFile(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	logPath := Config.LogPath

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		Log.FileLogger.Printf("ERROR: Log file not found: %v", err)
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "⛔ Лог-файл не найден."))
		return
	}

	fileBytes, err := os.ReadFile(logPath)
	if err != nil {
		Log.FileLogger.Printf("ERROR: Couldn't read log file: %v", err)
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "⛔ Ошибка при чтении лог-файла."))
		return
	}

	doc := tgbotapi.NewDocument(message.Chat.ID, tgbotapi.FileBytes{
		Name:  filepath.Base(logPath),
		Bytes: fileBytes,
	})
	_, err = bot.Send(doc)
	if err != nil {
		Log.FileLogger.Printf("ERROR: Couldn't send log file: %v", err)
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "⛔ Не удалось отправить лог-файл."))
	}
}

func (a *Admin) selectObject(bot *tgbotapi.BotAPI, message *tgbotapi.Message, updates tgbotapi.UpdatesChannel) string {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите объект:")
	msg.ReplyMarkup = a.generateKeyboard(Config.Objects, "Вернуться в главное меню")
	bot.Send(msg)

	for update := range updates {
		if update.Message != nil && update.Message.Text == "Вернуться в главное меню" {
			a.ShowAdminOptions(bot, update.Message)
			return ""
		}
		if update.Message != nil && handlers.IsObjectSelection(update.Message.Text) {
			return update.Message.Text
		}
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "⛔ Неверный выбор. Попробуйте снова."))
	}

	return ""
}

func (a *Admin) selectPosition(bot *tgbotapi.BotAPI, message *tgbotapi.Message, updates tgbotapi.UpdatesChannel) string {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Выберите должность:")
	msg.ReplyMarkup = a.generateKeyboard(Config.Positions, "Вернуться в главное меню")
	bot.Send(msg)

	for update := range updates {
		if update.Message != nil && update.Message.Text == "Вернуться в главное меню" {
			a.ShowAdminOptions(bot, update.Message)
			return ""
		}
		if update.Message != nil && handlers.IsPositionSelection(update.Message.Text) {
			return update.Message.Text
		}
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "⛔ Неверный выбор. Попробуйте снова."))
	}

	return ""
}

func (a *Admin) selectDate(bot *tgbotapi.BotAPI, message *tgbotapi.Message, updates tgbotapi.UpdatesChannel) string {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Введите дату в формате 'дд.мм.гггг':")
	bot.Send(msg)

	for update := range updates {
		if update.Message != nil && update.Message.Text == "Вернуться в главное меню" {
			a.ShowAdminOptions(bot, update.Message)
			return ""
		}
		if update.Message != nil && handlers.IsDateSelection(update.Message.Text) {
			return update.Message.Text
		}
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "⛔ Неверный формат даты. Попробуйте снова."))
	}

	return ""
}

func (a *Admin) sendPhotos(bot *tgbotapi.BotAPI, message *tgbotapi.Message, position, object, date string) {
	photosPath := filepath.Join(Config.BaseDir, position, object, date)
	entries, err := os.ReadDir(photosPath)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "⛔ Директория не найдена."))
		return
	}

	if len(entries) == 0 {
		bot.Send(tgbotapi.NewMessage(message.Chat.ID, "⛔ В выбранной папке нет фотографий."))
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			filePath := filepath.Join(photosPath, entry.Name())
			photoBytes, err := os.ReadFile(filePath)
			if err != nil {
				log.Printf("⛔ Ошибка при чтении файла %s: %v", filePath, err)
				continue
			}

			photo := tgbotapi.NewPhoto(message.Chat.ID, tgbotapi.FileBytes{Name: entry.Name(), Bytes: photoBytes})
			bot.Send(photo)
		}
	}
}

func (a *Admin) generateKeyboard(options []string, exitOption string) tgbotapi.ReplyKeyboardMarkup {
	var keyboardRows [][]tgbotapi.KeyboardButton
	row := []tgbotapi.KeyboardButton{}

	for i, option := range options {
		row = append(row, tgbotapi.NewKeyboardButton(option))
		if (i+1)%2 == 0 {
			keyboardRows = append(keyboardRows, row)
			row = []tgbotapi.KeyboardButton{}
		}
	}
	if len(row) > 0 {
		keyboardRows = append(keyboardRows, row)
	}
	keyboardRows = append(keyboardRows, []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton(exitOption),
	})

	return tgbotapi.ReplyKeyboardMarkup{
		Keyboard:        keyboardRows,
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}
}

func buildFileTree(path string, prefix string, depth int) string {
	entries, err := os.ReadDir(path)
	if err != nil {
		return "⛔ Ошибка при чтении директории: " + path
	}

	var tree strings.Builder
	for i, entry := range entries {
		isLast := i == len(entries)-1
		tree.WriteString(prefix)
		if isLast {
			tree.WriteString("└── ")
		} else {
			tree.WriteString("├── ")
		}
		tree.WriteString(entry.Name())
		if entry.IsDir() {
			tree.WriteString("/")
		}
		tree.WriteString("\n")

		if entry.IsDir() {
			newPrefix := prefix
			if isLast {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			tree.WriteString(buildFileTree(filepath.Join(path, entry.Name()), newPrefix, depth+1))
		}
	}
	return tree.String()
}
