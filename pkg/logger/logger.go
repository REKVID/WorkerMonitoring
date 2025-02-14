package Log

import (
	"log"
	"os"
)

// Экспортируемый логгер
var FileLogger *log.Logger

func InitLogger() {
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	//defer file.Close()

	FileLogger = log.New(file, "", log.LstdFlags)
	file.Sync()
	FileLogger.Println("Logger Created") // Тестовый вывод
}
