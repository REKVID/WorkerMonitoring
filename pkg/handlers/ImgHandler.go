package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func DownloadPhoto(bot *tgbotapi.BotAPI, message *tgbotapi.Message, position, objectName string) (string, error) {
	var fileID string

	// Проверяем, есть ли фотография в сообщении
	if len(message.Photo) > 0 {
		photo := message.Photo[len(message.Photo)-1] // Берём фото с наибольшим размером
		fileID = photo.FileID
	} else {
		return "", fmt.Errorf("no valid photo found")
	}

	// Получаем URL файла
	file, err := bot.GetFileDirectURL(fileID)
	if err != nil {
		return "", fmt.Errorf("error getting file URL: %v", err)
	}

	// Скачиваем файл
	resp, err := http.Get(file)
	if err != nil {
		return "", fmt.Errorf("error downloading file: %v", err)
	}
	defer resp.Body.Close()

	// Генерируем имя файла на основе Telegram username
	fileName := fmt.Sprintf("%s.jpg", message.From.UserName)

	// Создаём директории с именем объекта, должности и текущей даты
	currentTime := time.Now()
	dateDir := currentTime.Format("02.01.2006")
	photosDir := filepath.Join("photos", position, objectName, dateDir)
	err = os.MkdirAll(photosDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("error creating directory: %v", err)
	}

	// Путь для сохранения файла
	filePath := filepath.Join(photosDir, fileName)

	// Проверяем, существует ли уже файл с таким именем
	// Если существует, добавляем суффикс (1), (2), ...
	i := 1
	baseName := fileName[:len(fileName)-4] // Берем имя без расширения
	for {
		// Проверяем, существует ли файл с таким именем
		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			// Если файл не существует, выходим из цикла
			break
		}

		// Если файл существует, формируем новое имя с суффиксом
		fileName = fmt.Sprintf("%s(%d).jpg", baseName, i)
		filePath = filepath.Join(photosDir, fileName)
		i++
	}

	// Создаём файл и сохраняем его
	f, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error creating file: %v", err)
	}
	_, err = io.Copy(f, resp.Body)
	f.Close()

	if err != nil {
		return "", fmt.Errorf("error saving file: %v", err)
	}

	return filePath, nil
}
