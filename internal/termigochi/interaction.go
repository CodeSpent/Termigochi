package termigochi

import (
	"fmt"
	"termigochi/internal/config"
)

func StartOnboarding(config *config.Config) error {
	fmt.Println("Welcome to Termigochi!")

	fmt.Print("Enter your player name: ")
	var playerName string
	_, err := fmt.Scanln(&playerName)
	if err != nil {
		return err
	}
	config.PlayerName = playerName
	config.IsFirstRun = false

	err = config.SaveConfig()
	if err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}
