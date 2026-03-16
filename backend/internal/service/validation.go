package service

import (
	"errors"
	"strings"
)

func ValidatePetInput(name, species string, age int, description string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name is required")
	}

	if strings.TrimSpace(description) == "" {
		return errors.New("description is required")
	}

	if age < 0 {
		return errors.New("age must be 0 or greater")
	}

	if !IsValidSpecies(species) {
		return errors.New("species must be CAT, DOG, or BIRD")
	}

	return nil
}

func IsValidSpecies(species string) bool {
	switch strings.TrimSpace(species) {
	case "CAT", "DOG", "BIRD":
		return true
	default:
		return false
	}
}
