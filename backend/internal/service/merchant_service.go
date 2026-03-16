package service

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"pet-store/backend/internal/repository"

	"github.com/google/uuid"
)

type MerchantService struct {
	DB    *sql.DB
	PetDB *repository.PetRepository
}

func NewMerchantService(db *sql.DB) *MerchantService {
	return &MerchantService{
		DB:    db,
		PetDB: repository.NewPetRepository(db),
	}
}

func (s *MerchantService) GetUnsoldPets(storeID string) ([]repository.PetRecord, error) {
	return s.PetDB.GetUnsoldPetsByStoreID(storeID)
}

func (s *MerchantService) GetSoldPets(storeID string, start, end time.Time) ([]repository.PetRecord, error) {
	return s.PetDB.GetSoldPetsByStoreIDAndRange(storeID, start, end)
}

func (s *MerchantService) CreatePet(
	ctx context.Context,
	storeID string,
	name string,
	species string,
	age int,
	pictureBase64 string,
	description string,
	imageDir string,
) (repository.PetRecord, error) {
	if err := ValidatePetInput(name, species, age, description); err != nil {
		return repository.PetRecord{}, err
	}

	pictureURL, err := saveBase64Image(imageDir, pictureBase64)
	if err != nil {
		return repository.PetRecord{}, errors.New("failed to save image")
	}

	var pet repository.PetRecord
	err = s.DB.QueryRowContext(ctx, `
		INSERT INTO pets (store_id, name, species, age, picture_path, description)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, store_id, name, species, age, picture_path, description, created_at, status, sold_at
	`, storeID, name, species, age, pictureURL, description).Scan(
		&pet.ID,
		&pet.StoreID,
		&pet.Name,
		&pet.Species,
		&pet.Age,
		&pet.PictureURL,
		&pet.Description,
		&pet.CreatedAt,
		&pet.Status,
		&pet.SoldAt,
	)
	if err != nil {
		return repository.PetRecord{}, err
	}

	return pet, nil
}

func (s *MerchantService) RemovePetIfAvailable(ctx context.Context, petID string, merchantStoreID string) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var status string
	var storeID string

	err = tx.QueryRowContext(ctx, `
		SELECT status, store_id
		FROM pets
		WHERE id = $1
		FOR UPDATE
	`, petID).Scan(&status, &storeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("pet not found")
		}
		return err
	}

	if storeID != merchantStoreID {
		return errors.New("forbidden: cannot remove pets from another store")
	}

	if status != "AVAILABLE" {
		return errors.New("pet is already sold and cannot be removed")
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM pets WHERE id = $1`, petID); err != nil {
		return err
	}

	return tx.Commit()
}

func saveBase64Image(imageDir, payload string) (string, error) {
	data := payload

	if strings.Contains(payload, ",") {
		parts := strings.SplitN(payload, ",", 2)
		data = parts[1]
	}

	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(imageDir, 0755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s.jpg", uuid.NewString())
	fullPath := filepath.Join(imageDir, filename)

	if err := os.WriteFile(fullPath, decoded, 0644); err != nil {
		return "", err
	}

	return "/media/pet-images/" + filename, nil
}
