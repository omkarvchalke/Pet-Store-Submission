package repository

import (
	"database/sql"

	"github.com/google/uuid"
)

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) CreateOrderTx(tx *sql.Tx, customerID, storeID string) (string, error) {
	orderID := uuid.NewString()

	_, err := tx.Exec(`
		INSERT INTO orders (id, customer_id, store_id, purchased_at)
		VALUES ($1, $2, $3, NOW())
	`, orderID, customerID, storeID)
	if err != nil {
		return "", err
	}

	return orderID, nil
}

func (r *OrderRepository) CreateOrderItemTx(tx *sql.Tx, orderID, petID string) error {
	_, err := tx.Exec(`
		INSERT INTO order_items (id, order_id, pet_id)
		VALUES ($1, $2, $3)
	`, uuid.NewString(), orderID, petID)

	return err
}

func (r *OrderRepository) MarkPetSoldTx(tx *sql.Tx, petID string) error {
	_, err := tx.Exec(`
		UPDATE pets
		SET status = 'SOLD', sold_at = NOW()
		WHERE id = $1
	`, petID)

	return err
}
