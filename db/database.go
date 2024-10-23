package db

import (
	"database/sql"
	"fmt"
	"go-store/types"
)

func GetAllProducts(conn *sql.DB) ([]types.Product, error) {
	rows, err := conn.Query("SELECT product_name, in_stock, image_name FROM product")
	if err != nil {
		return nil, err
	}
	defer rows.Close() //defer = dont run until the end of the function

	var products []types.Product

	for rows.Next() {
		var product types.Product
		rows.Scan(&product.Name, &product.QuantityInStock, &product.Image)
		products = append(products, product)
	}

	return products, nil
}

func AddOrder(conn *sql.DB, ord types.Order) (int64, error) {
	result, err := conn.Exec("INSERT INTO orders (product_id, customer_id, quantity, price) VALUES (?, ?, ?, ?)",
		ord.ProductID, ord.CustomerID, ord.Quantity, ord.Price)
	if err != nil {
		return 0, fmt.Errorf("addOrder: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addOrder: %v", err)
	}
	return id, nil
}

func ProductByName(conn *sql.DB, name string) ([]types.Product, error) {
	var products []types.Product

	rows, err := conn.Query("SELECT id, product_name, image_name, price, in_stock FROM product WHERE product_name = ?", name)
	if err != nil {
		return nil, fmt.Errorf("productByName %q: %v", name, err)
	}
	defer rows.Close()

	for rows.Next() {
		var p types.Product

		if err := rows.Scan(&p.ID, &p.Name, &p.Image, &p.Price, &p.QuantityInStock); err != nil {
			return nil, fmt.Errorf("productByName %q: %v", name, err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func ProductByID(conn *sql.DB, id int64) (types.Product, error) {
	var p types.Product
	row := conn.QueryRow("SELECT id, product_name, price, in_stock, image_name FROM product WHERE id = ?", id)

	err := row.Scan(&p.ID, &p.Name, &p.Price, &p.QuantityInStock, &p.Image)
	if err != nil {
		if err == sql.ErrNoRows {
			return p, fmt.Errorf("productByID %d: no such product", id)
		}
		return p, fmt.Errorf("productByID %d: %v", id, err)
	}

	return p, nil
}

func CustomerByLastName(conn *sql.DB, lastName string) ([]types.Customer, error) {
	var customers []types.Customer
	rows, err := conn.Query("SELECT * FROM customer WHERE last_name = ?", lastName)
	if err != nil {
		return nil, fmt.Errorf("customerByLastName %q: %v", lastName, err)
	}
	defer rows.Close()

	for rows.Next() {
		var c types.Customer
		if err := rows.Scan(&c.ID, &c.FirstName, &c.LastName, &c.Email); err != nil {
			return nil, fmt.Errorf("customerByLastName %q: %v", lastName, err)
		}
		customers = append(customers, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

func OrdersByCustomer(conn *sql.DB, customerID int) ([]types.Order, error) {
	var orders []types.Order
	rows, err := conn.Query("SELECT * FROM orders WHERE customer_id = ?", customerID)
	if err != nil {
		return nil, fmt.Errorf("ordersByCustomer %d: %v", customerID, err)
	}
	defer rows.Close()

	for rows.Next() {
		var o types.Order
		if err := rows.Scan(&o.ID, &o.ProductID, &o.CustomerID, &o.Quantity, &o.Price, &o.Tax, &o.Donation, &o.Timestamp); err != nil {
			return nil, fmt.Errorf("ordersByCustomer %d: %v", customerID, err)
		}
		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
