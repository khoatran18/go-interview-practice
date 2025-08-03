package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Product represents a product in the inventory system
type Product struct {
	ID       int64
	Name     string
	Price    float64
	Quantity int
	Category string
}

// ProductStore manages product operations
type ProductStore struct {
	db *sql.DB
}

// NewProductStore creates a new ProductStore with the given database connection
func NewProductStore(db *sql.DB) *ProductStore {
	return &ProductStore{db: db}
}

// InitDB sets up a new SQLite database and creates the products table
func InitDB(dbPath string) (*sql.DB, error) {
	// TODO: Open a SQLite database connection
	// TODO: Create the products table if it doesn't exist
	// The table should have columns: id, name, price, quantity, category
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS products (id INTEGER PRIMARY KEY, name TEXT, price REAL, quantity INTEGER, category TEXT)`,
	)

	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// CreateProduct adds a new product to the database
func (ps *ProductStore) CreateProduct(product *Product) error {
	// TODO: Insert the product into the database
	// TODO: Update the product.ID with the database-generated ID

	db := ps.db
	res, err := db.Exec(
		`INSERT INTO products (name, price, quantity, category) VALUES (?, ?, ?, ?)`,
		product.Name, product.Price, product.Quantity, product.Category,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	product.ID = id

	return nil
}

// GetProduct retrieves a product by ID
func (ps *ProductStore) GetProduct(id int64) (*Product, error) {
	// TODO: Query the database for a product with the given ID
	// TODO: Return a Product struct populated with the data or an error if not found
	db := ps.db
	result := db.QueryRow(
		`SELECT * FROM products WHERE id = ?`, id,
	)

	p := &Product{}
	err := result.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with id %d not found", id)
		}
		return nil, err
	}
	return p, nil
}

// UpdateProduct updates an existing product
func (ps *ProductStore) UpdateProduct(product *Product) error {
	// TODO: Update the product in the database
	// TODO: Return an error if the product doesn't exist

	res, err := ps.db.Exec(
		`UPDATE products SET name = ?, price = ?, quantity = ?, category = ?
		WHERE id = ?`,
		product.Name, product.Price, product.Quantity, product.Category, product.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no product with id %d found", product.ID)
	}

	return nil
}

// DeleteProduct removes a product by ID
func (ps *ProductStore) DeleteProduct(id int64) error {
	// TODO: Delete the product from the database
	// TODO: Return an error if the product doesn't exist

	res, err := ps.db.Exec(
		`DELETE FROM products WHERE id = ?`,
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no product with id %d found", id)
	}

	return nil
}

// ListProducts returns all products with optional filtering by category
func (ps *ProductStore) ListProducts(category string) ([]*Product, error) {
	// TODO: Query the database for products
	// TODO: If category is not empty, filter by category
	// TODO: Return a slice of Product pointers
	var rows *sql.Rows
	var err error
	
	if category == "" {
		rows, err = ps.db.Query("SELECT id, name, price, quantity, category FROM products")
	} else {
		rows, err = ps.db.Query(
			`SELECT * FROM products WHERE category = ?`,
			category,
		)
	}

	defer rows.Close()
	if err != nil {
		return nil, err
	}

	products := []*Product{}
	for rows.Next() {
		p := &Product{}
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Category)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// BatchUpdateInventory updates the quantity of multiple products in a single transaction
func (ps *ProductStore) BatchUpdateInventory(updates map[int64]int) error {
	// TODO: Start a transaction
	// TODO: For each product ID in the updates map, update its quantity
	// TODO: If any update fails, roll back the transaction
	// TODO: Otherwise, commit the transaction

	tx, err := ps.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare(`UPDATE products Set quantity = ? WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for id, quantity := range updates {
		res, errExec := stmt.Exec(quantity, id)
		if errExec != nil {
			err = errExec
			return err
		}

		rowsEffect, errRa := res.RowsAffected()
		if errRa != nil {
			err = errRa
			return err
		}

		if rowsEffect == 0 {
			return fmt.Errorf("product with id %d not found", id)
		}
	}

	err = tx.Commit()
	return err
}

func main() {
	// Optional: you can write code here to test your implementation
}
