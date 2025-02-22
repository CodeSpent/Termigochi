package termigochi

import (
	"fmt"
	"termigochi/internal/config"
	"termigochi/internal/logger"
	"termigochi/internal/models"
)

func FeedPet(pet *models.Pet, foodInput string) {
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
	pet.SaveState(config.DefaultPetStateFilePath)
	logger.EventLogger.Printf("Fed %s, Hunger: %d\n", food.Name, pet.Hunger)
}

func PlayWithPet(pet *models.Pet, toyInput string) {
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
	pet.SaveState(config.DefaultPetStateFilePath)
	logger.EventLogger.Printf("Played with %s, Happiness: %d\n", toy.Name, pet.Happiness)
}

func PrintStatus(pet *models.Pet) {
	fmt.Printf("Hunger: %d, Happiness: %d\n", pet.Hunger, pet.Happiness)
}

func ReportState(pet *models.Pet) {
	hungerColor := GetColor(pet.Hunger)
	happinessColor := GetColor(pet.Happiness)
	healthColor := GetColor(pet.Health)

	fmt.Println()
	fmt.Printf("%s", pet.Name)
	fmt.Println()
	fmt.Println()

	if pet.Hatched {
		fmt.Printf("%sHealth: %d%s\n", healthColor, pet.Health, ResetColor())
		fmt.Printf("%sHunger: %d%s\n", hungerColor, pet.Hunger, ResetColor())
		fmt.Printf("%sHappiness: %d%s\n", happinessColor, pet.Happiness, ResetColor())
	} else {
		fmt.Printf("Still incubating.. Check back later.")
	}

	fmt.Println()
}

func GetColor(value int) string {
	switch {
	case value < 30:
		return "\033[31m" // Red
	case value < 70:
		return "\033[33m" // Yellow
	default:
		return "\033[32m" // Green
	}
}

func ResetColor() string {
	return "\033[0m"
}
