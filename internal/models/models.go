package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

type Pet struct {
	Name      string `json:"name"`
	Hunger    int    `json:"hunger"`
	Happiness int    `json:"happiness"`
	Food      []Food `json:"food"`
	Toys      []Toy  `json:"toys"`
}

func NewPet(name string) *Pet {
	return &Pet{
		Name:      name,
		Hunger:    50,
		Happiness: 50,
		Food:      make([]Food, 0),
		Toys:      make([]Toy, 0),
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
	return ioutil.WriteFile(filename, data, 0644)
}

func LoadState(filename string) (*Pet, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return NewPet("Pet"), nil
		}
		return nil, err
	}
	var pet Pet
	if err := json.Unmarshal(data, &pet); err != nil {
		return nil, err
	}
	return &pet, nil
}

func (p *Pet) UpdateState() {
	for {
		time.Sleep(1 * time.Second) // Adjust the time interval as needed
		p.Hunger -= 1
		p.Happiness -= 1

		if p.Hunger < 0 {
			p.Hunger = 0
		}
		if p.Happiness < 0 {
			p.Happiness = 0
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
