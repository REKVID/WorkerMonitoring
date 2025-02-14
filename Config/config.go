package Config

import (
	"encoding/json"
	"io/ioutil"
)

var (
	Positions []string
	Objects   []string
	BotToken  string
	AdminKey  string
	BaseDir   string
	LogPath   string
)

type config struct {
	Positions []string `json:"positions"`
	Objects   []string `json:"objects"`
	BotToken  string   `json:"bot_token"`
	AdminKey  string   `json:"admin_key"`
	BaseDir   string   `json:"base_dir"`
	LogPath   string   `json:"log_path"`
}

func LoadConfig(filePath string) error {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var cfg config
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return err
	}

	Positions = cfg.Positions
	Objects = cfg.Objects
	BotToken = cfg.BotToken
	AdminKey = cfg.AdminKey
	BaseDir = cfg.BaseDir
	LogPath = cfg.LogPath

	return nil
}

func SaveConfig(filePath string) error {
	cfg := config{
		Positions: Positions,
		Objects:   Objects,
		BotToken:  BotToken,
		AdminKey:  AdminKey,
		BaseDir:   BaseDir,
		LogPath:   LogPath,
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, data, 0644)
}
