package models

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"termigochi/internal/config"
	"time"
)

type Pet struct {
	Name             string    `json:"name"`
	Hunger           int       `json:"hunger"`
	Happiness        int       `json:"happiness"`
	Health           int       `json:"health"`
	Food             []Food    `json:"food"`
	Toys             []Toy     `json:"toys"`
	CreatedDate      time.Time `json:"created_date"`
	HatchDate        time.Time `json:"hatch_date"`
	Hatched          bool      `json:"hatched"`
	HatchEventPlayed bool      `json:"hatch_event_played"`
}

func NewPet(name string) *Pet {
	return &Pet{
		Name:        name,
		Hunger:      50,
		Happiness:   50,
		Health:      100,
		Food:        make([]Food, 0),
		Toys:        make([]Toy, 0),
		HatchDate:   generateHatchDate(),
		Hatched:     false,
		CreatedDate: time.Now(),
	}
}

func (p *Pet) Feed(food Food) {
	p.Hunger += food.Nutrition
	if p.Hunger > 100 {
		p.Hunger = 100
	}
}

func (p *Pet) Play(toy Toy) {
	p.Happiness += toy.FunLevel
	if p.Happiness > 100 {
		p.Happiness = 100
	}
}

func (p *Pet) SaveState(filename string) error {
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 644)
}

func (p *Pet) Hatch() *Pet {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Congratulations! What's their name?")

	name, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	p.Name = name

	// Update HatchDate to actual for accuracy
	p.HatchDate = time.Now()
	p.Hatched = true

	p.SaveState(config.DefaultPetStateFilePath)
	return p
}

func generateHatchDate() time.Time {
	hatchDuration := time.Duration(rand.Intn(60000)) * time.Millisecond
	return time.Now().Add(hatchDuration)
}

func LoadPetFromStateFile(filename string) (*Pet, error) {
	f, err := os.Open(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			f, err = os.Create(filename)
			if err != nil {
				fmt.Println("Error creating file:", err)
				return nil, err
			}

			pet := NewPet("Egg")
			err = pet.SaveState(filename)
			if err != nil {
				fmt.Println("Error creating pet:", err)
			}
			return pet, nil
		}
	}

	data, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}

	if err != nil {
		return nil, err
	}

	var pet Pet
	if err := json.Unmarshal(data, &pet); err != nil {
		return nil, err
	}
	return &pet, nil
}

func (p *Pet) TickState() {
	for {
		time.Sleep(1 * time.Second)

		// Check if egg is hatched & hatch event has played
		if !p.Hatched {
			return
		}

		// Check if egg is hatched & hatch event SHOULD play
		if p.HatchDate.After(time.Now()) {
			p.Hatch()
		}

		if p.Hunger != 0 {
			p.Hunger -= 1
		}

		if p.Happiness != 0 {
			p.Happiness -= 1
		}
	}
}

type Food struct {
	Name      string `json:"name"`
	Nutrition int    `json:"nutrition"`
}

type Toy struct {
	Name     string `json:"name"`
	FunLevel int    `json:"fun_level"`
}

func NewFood(name string, nutrition int) Food {
	return Food{
		Name:      name,
		Nutrition: nutrition,
	}
}

func NewToy(name string, funLevel int) Toy {
	return Toy{
		Name:     name,
		FunLevel: funLevel,
	}
}
