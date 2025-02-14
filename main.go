package main

import (
	"TgTaskBot/Config"
	"TgTaskBot/Log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	Log.InitLogger()

	// Загрузка конфигурации
	err := Config.LoadConfig("./Config/configuration.json")
	if err != nil {
		Log.FileLogger.Printf("ERROR: Failed to load config: %v", err)
		os.Exit(1)
	}
	Log.FileLogger.Printf("INFO: Config loaded successfully.")

	// Инициализация бота
	botAPI, err := tgbotapi.NewBotAPI(Config.BotToken)
	if err != nil {
		Log.FileLogger.Printf("FATAL: Bot initialization failed: %v", err)
		os.Exit(1)
	}
	Log.FileLogger.Printf("INFO: Authorized as %s", botAPI.Self.UserName)

	// Настройка канала обновлений
	updates := botAPI.GetUpdatesChan(tgbotapi.NewUpdate(0))
	Log.FileLogger.Printf("INFO: Starting updates processing...")

	// Обработка входящих сообщений
	for update := range updates {
		if update.Message != nil {
			//Log.FileLogger.Printf(
			//	"DEBUG: Received message from @%s (ID:%d): %s",
			//	update.Message.From.UserName,
			//	update.Message.From.ID,
			//	update.Message.Text,
			//)

			HandleMessage(botAPI, update.Message, updates)
		}
	}
}
