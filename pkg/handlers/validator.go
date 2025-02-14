package handlers

import (
	"TgTaskBot/Config"
	"time"
)

func IsDateSelection(text string) bool {
	_, err := time.Parse("02.01.2006", text)
	return err == nil
}

func IsObjectSelection(text string) bool {
	for _, object := range Config.Objects {
		if text == object {
			return true
		}
	}
	return false
}

func IsPositionSelection(text string) bool {
	for _, position := range Config.Positions {
		if text == position {
			return true
		}
	}
	return false
}
