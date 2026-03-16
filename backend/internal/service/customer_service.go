package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"pet-store/backend/internal/repository"

	"github.com/lib/pq"
)

type CustomerService struct {
	DB      *sql.DB
	PetDB   *repository.PetRepository
	OrderDB *repository.OrderRepository
}

func NewCustomerService(db *sql.DB) *CustomerService {
	return &CustomerService{
		DB:      db,
		PetDB:   repository.NewPetRepository(db),
		OrderDB: repository.NewOrderRepository(db),
	}
}

func (s *CustomerService) GetAvailablePets(storeSlug string) ([]repository.PetRecord, error) {
	return s.PetDB.GetAvailablePetsByStoreSlug(storeSlug)
}

func (s *CustomerService) PurchasePet(ctx context.Context, customerID, petID string) (bool, string, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return false, "", err
	}
	defer tx.Rollback()

	var name, status, storeID string

	err = tx.QueryRowContext(ctx, `
		SELECT name, status, store_id
		FROM pets
		WHERE id = $1
		FOR UPDATE
	`, petID).Scan(&name, &status, &storeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, "Pet not found", nil
		}
		return false, "", err
	}

	if status != "AVAILABLE" {
		return false, fmt.Sprintf("Sorry, %s is no longer available.", name), nil
	}

	orderID, err := s.OrderDB.CreateOrderTx(tx, customerID, storeID)
	if err != nil {
		return false, "", err
	}

	if err := s.OrderDB.MarkPetSoldTx(tx, petID); err != nil {
		return false, "", err
	}

	if err := s.OrderDB.CreateOrderItemTx(tx, orderID, petID); err != nil {
		return false, "", err
	}

	if err := tx.Commit(); err != nil {
		return false, "", err
	}

	return true, fmt.Sprintf("You successfully purchased %s.", name), nil
}

func (s *CustomerService) CheckoutCart(ctx context.Context, customerID string, petIDs []string) (bool, string, []string, error) {
	if len(petIDs) == 0 {
		return false, "Your cart is empty.", []string{}, nil
	}

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return false, "", nil, err
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, `
		SELECT id, name, status, store_id
		FROM pets
		WHERE id = ANY($1)
		FOR UPDATE
	`, pq.Array(petIDs))
	if err != nil {
		return false, "", nil, err
	}
	defer rows.Close()

	var pets []repository.CartPetRecord
	var unavailable []string
	storeSet := map[string]bool{}

	for rows.Next() {
		var row repository.CartPetRecord
		if err := rows.Scan(&row.ID, &row.Name, &row.Status, &row.StoreID); err != nil {
			return false, "", nil, err
		}
		pets = append(pets, row)
		storeSet[row.StoreID] = true

		if row.Status != "AVAILABLE" {
			unavailable = append(unavailable, row.Name)
		}
	}

	if len(pets) != len(petIDs) {
		return false, "Some pets could not be found.", []string{}, nil
	}

	if len(storeSet) != 1 {
		return false, "All pets in checkout must belong to the same store.", []string{}, nil
	}

	if len(unavailable) > 0 {
		return false, fmt.Sprintf("Some pets are no longer available: %v", unavailable), unavailable, nil
	}

	storeID := pets[0].StoreID
	orderID, err := s.OrderDB.CreateOrderTx(tx, customerID, storeID)
	if err != nil {
		return false, "", nil, err
	}

	for _, pet := range pets {
		if err := s.OrderDB.MarkPetSoldTx(tx, pet.ID); err != nil {
			return false, "", nil, err
		}

		if err := s.OrderDB.CreateOrderItemTx(tx, orderID, pet.ID); err != nil {
			return false, "", nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return false, "", nil, err
	}

	return true, "Checkout successful.", []string{}, nil
}
