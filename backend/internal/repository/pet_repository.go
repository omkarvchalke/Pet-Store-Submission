package repository

import (
	"database/sql"
	"time"
)

type PetRecord struct {
	ID          string
	StoreID     string
	Name        string
	Species     string
	Age         int
	PictureURL  string
	Description string
	CreatedAt   time.Time
	Status      string
	SoldAt      sql.NullTime
}

type CartPetRecord struct {
	ID      string
	Name    string
	Status  string
	StoreID string
}

type PetRepository struct {
	DB *sql.DB
}

func NewPetRepository(db *sql.DB) *PetRepository {
	return &PetRepository{DB: db}
}

func (r *PetRepository) GetAvailablePetsByStoreSlug(storeSlug string) ([]PetRecord, error) {
	rows, err := r.DB.Query(`
		SELECT p.id, p.store_id, p.name, p.species, p.age, p.picture_path, p.description, p.created_at, p.status, p.sold_at
		FROM pets p
		JOIN stores s ON s.id = p.store_id
		WHERE s.slug = $1 AND p.status = 'AVAILABLE'
		ORDER BY p.created_at DESC
	`, storeSlug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pets []PetRecord
	for rows.Next() {
		var pet PetRecord
		if err := rows.Scan(
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
		); err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}

	return pets, nil
}

func (r *PetRepository) GetUnsoldPetsByStoreID(storeID string) ([]PetRecord, error) {
	rows, err := r.DB.Query(`
		SELECT id, store_id, name, species, age, picture_path, description, created_at, status, sold_at
		FROM pets
		WHERE store_id = $1 AND status = 'AVAILABLE'
		ORDER BY created_at DESC
	`, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pets []PetRecord
	for rows.Next() {
		var pet PetRecord
		if err := rows.Scan(
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
		); err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}

	return pets, nil
}

func (r *PetRepository) GetSoldPetsByStoreIDAndRange(storeID string, start, end time.Time) ([]PetRecord, error) {
	rows, err := r.DB.Query(`
		SELECT id, store_id, name, species, age, picture_path, description, created_at, status, sold_at
		FROM pets
		WHERE store_id = $1
		  AND status = 'SOLD'
		  AND sold_at >= $2
		  AND sold_at <= $3
		ORDER BY sold_at DESC
	`, storeID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pets []PetRecord
	for rows.Next() {
		var pet PetRecord
		if err := rows.Scan(
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
		); err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}

	return pets, nil
}
