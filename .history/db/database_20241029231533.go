package db

import (
	"database/sql"
	"fmt"
	"go-store/types"
	"time"
)

// CUSTOMER FUNCTIONS //

func GetAllCustomers(conn *sql.DB) ([]types.Customer, error) {
	stmt, err := conn.Prepare("SELECT id, first_name, last_name, email FROM customer")
	if err != nil {
		return nil, fmt.Errorf("GetAllCustomers: prepare: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("GetAllCustomers: query: %v", err)
	}
	defer rows.Close()

	var customers []types.Customer
	for rows.Next() {
		var c types.Customer
		if err := rows.Scan(&c.ID, &c.FirstName, &c.LastName, &c.Email); err != nil {
			return nil, fmt.Errorf("GetAllCustomers: scan: %v", err)
		}
		customers = append(customers, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

func CustomerByID(conn *sql.DB, id int64) (types.Customer, error) {
	stmt, err := conn.Prepare("SELECT id, first_name, last_name, email FROM customer WHERE id = ?")
	if err != nil {
		return types.Customer{}, fmt.Errorf("CustomerByID: prepare: %v", err)
	}
	defer stmt.Close()

	var c types.Customer
	err = stmt.QueryRow(id).Scan(&c.ID, &c.FirstName, &c.LastName, &c.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c, fmt.Errorf("CustomerByID %d: no such customer", id)
		}
		return c, fmt.Errorf("CustomerByID %d: %v", id, err)
	}

	return c, nil
}

func CustomerByLastName(conn *sql.DB, lastName string) ([]types.Customer, error) {
	// Prepare statement
	stmt, err := conn.Prepare("SELECT * FROM customer WHERE last_name = ?")
	if err != nil {
		return nil, fmt.Errorf("CustomerByLastName: prepare: %v", err)
	}
	defer stmt.Close()

	// Execute
	rows, err := stmt.Query(lastName)
	if err != nil {
		return nil, fmt.Errorf("CustomerByLastName: query: %v", err)
	}
	defer rows.Close()

	// Store results
	var customers []types.Customer
	for rows.Next() {
		var c types.Customer

		if err := rows.Scan(&c.ID, &c.FirstName, &c.LastName, &c.Email); err != nil {
			return nil, fmt.Errorf("CustomerByLastName: scan: %v", err)
		}
		customers = append(customers, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

func CustomerByEmail(conn *sql.DB, email string) (types.Customer, error) {
	stmt, err := conn.Prepare("SELECT id, first_name, last_name, email FROM customer WHERE email = ?")
	if err != nil {
		return types.Customer{}, fmt.Errorf("CustomerByEmail: prepare: %v", err)
	}
	defer stmt.Close()

	var c types.Customer
	err = stmt.QueryRow(email).Scan(&c.ID, &c.FirstName, &c.LastName, &c.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c, sql.ErrNoRows
		}
		return c, fmt.Errorf("CustomerByEmail: %v", err)
	}

	return c, nil
}

func AddCustomer(conn *sql.DB, firstName, lastName, email string) (int64, error) {
	stmt, err := conn.Prepare("INSERT INTO customer (first_name, last_name, email) VALUES (?, ?, ?)")
	if err != nil {
		return 0, fmt.Errorf("AddCustomer: prepare: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(firstName, lastName, email)
	if err != nil {
		return 0, fmt.Errorf("AddCustomer: exec: %v", err)
	}

	customerID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("AddCustomer: %v", err)
	}

	return customerID, nil
}

// PRODUCT FUNCTIONS //

func GetAllProducts(conn *sql.DB) ([]types.Product, error) {
	stmt, err := conn.Prepare("SELECT id, product_name, image_name, price, in_stock FROM product")
	if err != nil {
		return nil, fmt.Errorf("GetAllProducts: prepare: %v", err)
	}
	defer stmt.Close() //defer = dont run until the end of the function

	rows, err := stmt.Query() // run the query + error check
	if err != nil {
		return nil, fmt.Errorf("GetAllProducts: query: %v", err)
	}
	defer rows.Close()

	var products []types.Product

	for rows.Next() {
		var product types.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Image, &product.Price, &product.QuantityInStock); err != nil {
			return nil, fmt.Errorf("GetAllProducts: scan: %v", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func ProductByName(conn *sql.DB, name string) ([]types.Product, error) {
	// Prepare statement
	stmt, err := conn.Prepare("SELECT id, product_name, image_name, price, in_stock FROM product WHERE product_name = ?")
	if err != nil {
		return nil, fmt.Errorf("ProductByName: prepare: %v", err)
	}
	defer stmt.Close()

	// Execute
	rows, err := stmt.Query(name)
	if err != nil {
		return nil, fmt.Errorf("ProductByName: query: %v", err)
	}
	defer rows.Close()

	//store results
	var products []types.Product
	for rows.Next() {
		var p types.Product

		if err := rows.Scan(&p.ID, &p.Name, &p.Image, &p.Price, &p.QuantityInStock); err != nil {
			return nil, fmt.Errorf("ProductByName: scan: %v", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func ProductByID(conn *sql.DB, id int64) (types.Product, error) {
	// Prepare statement
	stmt, err := conn.Prepare("SELECT id, product_name, price, in_stock, image_name FROM product WHERE id = ?")
	if err != nil {
		return types.Product{}, fmt.Errorf("ProductByID: prepare: %v", err)
	}
	defer stmt.Close()

	// Execute
	var p types.Product
	err = stmt.QueryRow(id).Scan(&p.ID, &p.Name, &p.Price, &p.QuantityInStock, &p.Image)

	if err != nil {
		if err == sql.ErrNoRows {
			return p, fmt.Errorf("ProductByID %d: no such product", id)
		}
		return p, fmt.Errorf("ProductByID %d: %v", id, err)
	}

	return p, nil
}

func UpdateProductQuantity(conn *sql.DB, productID int64, quantitySold int) error {
	// Use a single update query to avoid concurrency issues
	stmt, err := conn.Prepare(`UPDATE product SET in_stock = in_stock - ? WHERE id = ? AND in_stock >= ?`)
	if err != nil {
		return fmt.Errorf("UpdateProductQuantity: prepare update: %v", err)
	}
	defer stmt.Close()

	// Execute the update query
	res, err := stmt.Exec(quantitySold, productID, quantitySold)
	if err != nil {
		return fmt.Errorf("UpdateProductQuantity: exec: %v", err)
	}

	// Check if rows were affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("UpdateProductQuantity: rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("UpdateProductQuantity: no rows affected, product ID may be incorrect or insufficient stock")
	}

	fmt.Printf("Selling %d units of product ID %d\n", quantitySold, productID)

	return nil
}

// ORDER FUNCTIONS //

func GetAllOrders(conn *sql.DB) ([]types.Order, error) {
    stmt, err := conn.Prepare(`
        SELECT 
            o.customer_first, o.customer_last, o.product_name, 
            o.quantity, o.price, IFNULL(o.tax, 0), IFNULL(o.donation, 0), o.timestamp
        FROM orders o
        JOIN customer c ON c.first_name = o.customer_first AND c.last_name = o.customer_last
        JOIN product p ON p.product_name = o.product_name
    `)

	if err != nil {
		return nil, fmt.Errorf("GetAllOrders: prepare: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("GetAllOrders: query: %v", err)
	}
	defer rows.Close()

	var orders []types.Order
	for rows.Next() {
        var o types.OrderDetails
        // Updated to include timestamp
        if err := rows.Scan(&o.CustomerFirstName, &o.CustomerLastName, &o.ProductName, &o.Quantity, &o.Price, &o.Tax, &o.Donation, &o.Timestamp); err != nil {
            return nil, fmt.Errorf("GetAllOrders: scan: %v", err)
        }
        // Convert timestamp to a human-readable string for display
        o.ReadableTimestamp = time.UnixMilli(o.Timestamp).Format("2006-01-02 15:04:05 MST")
        orders = append(orders, o)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return orders, nil
}

func AddOrder(conn *sql.DB, ord types.Order) (int64, error) {
	stmt, err := conn.Prepare(`
        INSERT INTO orders (customer_first, customer_last, product_name, quantity, price, tax, donation, timestamp)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("AddOrder: prepare: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		ord.CustomerFirstName, ord.CustomerLastName, ord.ProductName, ord.Quantity, ord.Price, ord.Tax, ord.Donation, ord.Timestamp)
	if err != nil {
		return 0, fmt.Errorf("AddOrder: exec: %v", err)
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("AddOrder: %v", err)
	}

	return orderID, nil
}

func OrdersByCustomer(conn *sql.DB, customerFirstName string, customerLastName string) ([]types.Order, error) {

	stmt, err := conn.Prepare(`
        SELECT customer_first, customer_last, product_name, quantity, price, IFNULL(tax, 0), IFNULL(donation, 0)
        FROM orders
        WHERE customer_first = ? AND customer_last = ?`)
	if err != nil {
		return nil, fmt.Errorf("OrdersByCustomer: prepare: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(customerFirstName, customerLastName)
	if err != nil {
		return nil, fmt.Errorf("OrdersByCustomer: query: %v", err)
	}
	defer rows.Close()

	var orders []types.Order
	for rows.Next() {
		var o types.Order

		if err := rows.Scan(&o.CustomerFirstName, &o.CustomerLastName, &o.ProductName, &o.Quantity, &o.Price, &o.Tax, &o.Donation); err != nil {
			return nil, fmt.Errorf("OrdersByCustomer: scan: %v", err)
		}
		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func OrderExists(conn *sql.DB, timestamp int64) (bool, error) {
	stmt, err := conn.Prepare("SELECT COUNT(*) FROM orders WHERE timestamp = ?")
	if err != nil {
		return false, fmt.Errorf("OrderExists: prepare: %v", err)
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(timestamp).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("OrderExists: query: %v", err)
	}

	return count > 0, nil
}
