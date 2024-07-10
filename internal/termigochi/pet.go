package termigochi

import (
	"fmt"
	"os"
	"termigochi/internal/logger"
	"termigochi/internal/models"
	"time"
)

func petBackgroundService() {

	pet, err := models.LoadState(stateFile)
	if err != nil {
		logger.ServiceLogger.Printf("Error loading state: %v\n", err)
		os.Exit(1)
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		updatePetState(pet)
		saveAndLogState(pet)
	}
}

func updatePetState(pet *models.Pet) {
	pet.Hunger--
	pet.Happiness--

	if pet.Hunger < 0 {
		pet.Hunger = 0
	}
	if pet.Happiness < 0 {
		pet.Happiness = 0
	}
}

func saveAndLogState(pet *models.Pet) {
	SaveState(pet)
	logger.EventLogger.Printf("Updated Hunger: %d, Happiness: %d\n", pet.Hunger, pet.Happiness)
	logger.ServiceLogger.Println("Pet state updated.")
}

func feedPet(pet *models.Pet, foodInput string) {
	var food models.Food
	switch foodInput {
	case "apple":
		food = models.NewFood("Apple", 20)
	case "banana":
		food = models.NewFood("Banana", 15)
	default:
		fmt.Println("Unknown food.")
		return
	}
	pet.Feed(food)
	fmt.Printf("You feed the Pet a %s!\n", food.Name)
	SaveState(pet)
	logger.EventLogger.Printf("Fed %s, Hunger: %d\n", food.Name, pet.Hunger)
}

func playWithPet(pet *models.Pet, toyInput string) {
	var toy models.Toy
	switch toyInput {
	case "ball":
		toy = models.NewToy("Ball", 20)
	case "doll":
		toy = models.NewToy("Doll", 15)
	default:
		fmt.Println("Unknown toy.")
		return
	}
	pet.Play(toy)
	fmt.Printf("You play with the Pet using a %s!\n", toy.Name)
	SaveState(pet)
	logger.EventLogger.Printf("Played with %s, Happiness: %d\n", toy.Name, pet.Happiness)
}

func printStatus(pet *models.Pet) {
	fmt.Printf("Hunger: %d, Happiness: %d\n", pet.Hunger, pet.Happiness)
}

func reportState(pet *models.Pet) {
	hungerColor := GetColor(pet.Hunger)
	happinessColor := GetColor(pet.Happiness)

	fmt.Println()
	fmt.Printf("%sHunger: %d%s\n", hungerColor, pet.Hunger, ResetColor())
	fmt.Printf("%sHappiness: %d%s\n", happinessColor, pet.Happiness, ResetColor())
	fmt.Println()
}
