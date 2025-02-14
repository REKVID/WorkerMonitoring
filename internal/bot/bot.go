package bot

import (
	"TgTaskBot/Config"
	Admin "TgTaskBot/internal/admin"
	Worker "TgTaskBot/internal/worker"
	Log "TgTaskBot/pkg/logger"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserState struct {
	Role   string
	Worker *Worker.Worker
}

var userStates sync.Map

// getUserState возвращает состояние пользователя или создает новое, если оно отсутствует
func getUserState(userID int64) *UserState {
	state, exists := userStates.Load(userID)
	if !exists {
		newState := &UserState{
			Worker: &Worker.Worker{},
		}
		userStates.Store(userID, newState)
		return newState
	}
	return state.(*UserState)
}

// ShowRoleSelection отображает клавиатуру для выбора роли
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

// HandleMessage обрабатывает входящие сообщения
func HandleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, updates tgbotapi.UpdatesChannel) {
	userID := message.From.ID
	userState := getUserState(userID)

	// Логирование входящего сообщения
	Log.FileLogger.Printf(
		"DEBUG: Received message from @%s: %s",
		message.From.UserName,
		message.Text,
	)

	// Обработка фото, если оно есть
	if message.Photo != nil {
		Log.FileLogger.Printf("INFO: User %s sent a photo", message.From.UserName)
		userState.Worker.HandlePhoto(bot, message)
		return
	}

	// Общая обработка команд
	commandHandlers := map[string]func(){
		"/start": func() {
			Log.FileLogger.Printf("INFO: User %s started the bot", message.From.UserName)
			ShowRoleSelection(bot, message)
		},
		"Вернуться в главное меню": func() {
			Log.FileLogger.Printf("INFO: User %s returned to the main menu", message.From.UserName)
			ShowRoleSelection(bot, message)
		},
		"Работник": func() {
			userState.Role = "Работник"
			Log.FileLogger.Printf("INFO: User %s selected role: Worker", message.From.UserName)
			userState.Worker.StartPhotoUploadProcess(bot, message, updates)
		},
		"Администратор": func() {
			userState.Role = "Администратор"
			Log.FileLogger.Printf("INFO: User %s selected role: Admin", message.From.UserName)
			Admin.SelectedAdmin.RequestAdminKey(bot, message, updates)
		},
		"Добавить объект": func() {
			Log.FileLogger.Printf("INFO: User %s is adding an object", message.From.UserName)
			Admin.SelectedAdmin.AddObject(bot, message, updates)
		},
		"Добавить должность": func() {
			Log.FileLogger.Printf("INFO: User %s is adding a position", message.From.UserName)
			Admin.SelectedAdmin.AddPosition(bot, message, updates)
		},
		"Просмотр фотографий": func() {
			Log.FileLogger.Printf("INFO: User %s is viewing photos", message.From.UserName)
			Admin.SelectedAdmin.GetFilteredPhotos(bot, message, updates)
		},
		"Просмотр файловой системы": func() {
			Log.FileLogger.Printf("INFO: User %s is viewing the file system", message.From.UserName)
			Admin.SelectedAdmin.ListFileSystem(bot, message, Config.BaseDir)
		},
		"просмотр файла логирования": func() {
			Log.FileLogger.Printf("INFO: User %s is viewing the log file", message.From.UserName)
			Admin.SelectedAdmin.SendLogFile(bot, message)
		},
	}

	// Если это команда, обрабатываем её
	if handler, exists := commandHandlers[message.Text]; exists {
		handler()
		return
	}

	// Логирование неизвестной команды
	Log.FileLogger.Printf("WARNING: Unknown command from user %s: %s", message.From.UserName, message.Text)
	msg := tgbotapi.NewMessage(message.Chat.ID, "Неизвестная команда. Пожалуйста, выберите действие из меню или используйте /start.")
	if _, err := bot.Send(msg); err != nil {
		Log.FileLogger.Printf("ERROR: Failed to send unknown command message: %v", err)
	}
}
