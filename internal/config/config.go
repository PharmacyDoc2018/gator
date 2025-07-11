package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	err := write(c)
	if err != nil {
		return err
	}
	return nil
}

func write(c *Config) error {
	data, err := json.Marshal(*c)
	if err != nil {
		return err
	}

	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	configFile, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	_, err = configFile.Write(data)
	if err != nil {
		return err
	}
	err = configFile.Close()
	if err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + "/" + configFileName, nil

}

func Read() (*Config, error) {
	gatorConfig := &Config{}

	configPath, err := getConfigFilePath()
	if err != nil {
		return gatorConfig, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return gatorConfig, err
	}

	err = json.Unmarshal(data, gatorConfig)
	if err != nil {
		return gatorConfig, err
	}

	return gatorConfig, nil

}
