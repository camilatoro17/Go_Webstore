package main

import (
	"net/http"
	"strconv"

	"database/sql"
	"fmt"
	"log"

	"github.com/a-h/templ"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	"go-store/db"
	"go-store/templates"
	"go-store/types"
)

// var products = map[string]float64{
// 	"2024 G80 M3":   78000,
// 	"2024 S63 AMG":  183000,
// 	"2024 Audi RS7": 128000,
// }

// var productIDs = map[string]int{
// 	"2024 G80 M3":   1,
// 	"2024 S63 AMG":  2,
// 	"2024 Audi RS7": 3,
// }

var conn *sql.DB

func main() {
	e := echo.New()

	// connection
	cfg := mysql.Config{
		User:   "ctorosuarez",
		Passwd: "skies",
		DBName: "ctorosuarez",
	}

	var err error
	conn, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		e.Logger.Fatal(err)
	}

	// e.Logger.Fatal(conn.Ping())

	// TODO: Fill in your products here with name -> price as the key -> value pair.
	// products := map[string]float64{
	// }
	//e.Use(etag.Etag()) //MAYBE UNCOMMENT THIS???? IDK???

	// INFO: If you wanted to load a CSS file, you'd do something like this:
	// `<link rel="stylesheet" href="assets/styles/styles.css">`
	e.Static("assets", "./assets")

	// TODO: Render your base store page here
	e.GET("/store", func(ctx echo.Context) error {

		// Get products from the db
		products, err := db.GetAllProducts(conn)
		if err != nil {
			e.Logger.Errorf("Error fetching products: %v", err)
			return ctx.String(http.StatusInternalServerError, "Error loading products")
		}

		// Filter out inactive products
		var activeProducts []types.Product
		for _, product := range products {
			if product.Inactive == 0 {
				activeProducts = append(activeProducts, product)
			}
		}

		// Render the store page with the products
		return Render(ctx, http.StatusOK, templates.Base(templates.Store(activeProducts)))
	})

	e.GET("/dbQueries", func(ctx echo.Context) error {
		// get all customers
		customers, err := db.GetAllCustomers(conn)
		if err != nil {
			e.Logger.Errorf("%+v", err)
			return ctx.String(http.StatusInternalServerError, "Error retrieving customers")
		}

		// output messages
		outputmsg := ""

		// customer by ID
		customerByID, err := db.CustomerByID(conn, 2)
		if err != nil {
			outputmsg += fmt.Sprintf("Customer 2 by id... Error: %+v\n", err)
		} else {
			outputmsg += fmt.Sprintf("Customer 2 by id... %s\n", customerByID.Email)
		}

		// customer by ID - DNE
		_, err = db.CustomerByID(conn, 3)
		if err != nil {
			outputmsg += "Customer 3? Customer 3 not found!\n"
		}

		// customer by email
		customerByEmail, err := db.CustomerByEmail(conn, "mmouse@mines.edu")
		if err != nil {
			outputmsg += fmt.Sprintf("Customer by email: Error: %+v\n", err)
		} else {
			outputmsg += fmt.Sprintf("Customer by email: %s\n", customerByEmail.Email)
		}

		// customer DNE
		randomEmail := "random@email.com"
		_, err = db.CustomerByEmail(conn, randomEmail)
		if err != nil {
			outputmsg += fmt.Sprintf("Customer by email exists? Customer %s not found... adding...\n", randomEmail)

			// Add the new customer
			customerID, err := db.AddCustomer(conn, "Not", "Found", randomEmail)
			if err != nil {
				e.Logger.Errorf("Error adding customer: %+v", err)
			} else {
				outputmsg += fmt.Sprintf("Added customer with ID: %d\n", customerID)
			}
		}

		// get all orders
		orders, err := db.GetAllOrders(conn)
		if err != nil {
			e.Logger.Errorf("%+v", err)
			return ctx.String(http.StatusInternalServerError, "Error retrieving orders")
		}

		// get products
		products, err := db.GetAllProducts(conn)
		if err != nil {
			e.Logger.Errorf("%+v", err)
			return ctx.String(http.StatusInternalServerError, "Error retrieving products")
		}

		return Render(ctx, http.StatusOK, templates.Queries(customers, len(customers), orders, len(orders), products, outputmsg))
	})

	// TODO: Handle the form submission and return the purchase confirmation view
	e.POST("/purchase", func(ctx echo.Context) error {
		// Grab the form details from ctx.FormValue("...")

		fname := ctx.FormValue("fname")
		lname := ctx.FormValue("lname")
		email := ctx.FormValue("email")
		carIDStr := ctx.FormValue("car")
		quantityStr := ctx.FormValue("quantity")
		roundup := ctx.FormValue("donate")
		timestampStr := ctx.FormValue("timestamp")

		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			fmt.Printf("Invalid timestamp: %v\n", err)
			return ctx.String(http.StatusBadRequest, "Invalid timestamp")
		}

		// check if order with the same timestamp exists
		exists, err := db.OrderExists(conn, timestamp)
		if err != nil {
			log.Printf("Error checking order existence: %v", err)
			return ctx.String(http.StatusInternalServerError, "Error checking order existence")
		}

		if exists {
			fmt.Println("Duplicate order detected. Skipping insertion.")
			return ctx.String(http.StatusOK, "This order has already been processed.")
		}

		// convert car ID from string to int
		carID, err := strconv.ParseInt(carIDStr, 10, 64)
		if err != nil {
			fmt.Printf("Invalid car ID: %v\n", err)
			return ctx.String(http.StatusBadRequest, "Invalid car ID")
		}

		// get product details from the db using car ID
		product, err := db.ProductByID(conn, carID)
		if err != nil {
			fmt.Printf("ProductByID error: %v\n", err)
			return ctx.String(http.StatusNotFound, "Car not found")
		}

		// convert quantity from str to int
		quantity, err := strconv.Atoi(quantityStr)
		if err != nil {
			fmt.Printf("Invalid quantity: %v\n", err)
			return ctx.String(http.StatusBadRequest, "Invalid quantity")
		}

		isNewCustomer := false

		// check if returning customer
		_, err = db.CustomerByEmail(conn, email)
		if err != nil {
			if err == sql.ErrNoRows {

				// new customer
				_, err = db.AddCustomer(conn, fname, lname, email)
				if err != nil {
					log.Printf("Error adding customer: %v", err)
					return ctx.String(http.StatusInternalServerError, "Error adding new customer")
				}
				log.Println("Welcome, new customer!")
				isNewCustomer = true
			} else {
				log.Printf("Unexpected error finding customer: %v", err)
				return ctx.String(http.StatusInternalServerError, "Unexpected error occurred while finding the customer")
			}
		} else {
			// prev customer
			log.Println("Welcome back, valued customer!")
		}

		// product, err = db.ProductByID(conn, carID)
		// if err != nil {
		// 	fmt.Printf("Error fetching product: %v", err)
		// 	return ctx.String(http.StatusInternalServerError, "Error fetching product")
		// }

		// check stock
		if product.QuantityInStock < quantity {
			return ctx.String(http.StatusBadRequest, "Not enough stock available")
		}

		// update product quantity after purchase
		err = db.UpdateProductQuantity(conn, product.ID, quantity)
		if err != nil {
			fmt.Printf("Error updating product quantity: %v\n", err)
			return ctx.String(http.StatusInternalServerError, "Error updating product quantity")
		}

		// Calculate subtotal (price * quantity)
		subtotal := product.Price * float64(quantity)

		// Calculate tax (2.9%)
		const taxRate = 0.029
		tax := subtotal * taxRate

		// Calculate total including tax
		totalWithTax := subtotal + tax

		// Check if user opted to round up
		var grandtotal float64
		var donation float64 = 0

		if roundup == "yes" {
			grandtotal = totalWithTax + 1.00 // Round up
			donation = 1.00
		} else {
			grandtotal = totalWithTax
		}

		// Create a personalized welcome message
		var welcomeMessage string
		if isNewCustomer {
			welcomeMessage = fmt.Sprintf("Thank you for your order, %s! Welcome new customer.", fname)
		} else {
			welcomeMessage = fmt.Sprintf("Welcome back, %s! Thank you for your order.", fname)
		}

		// Create the purchaseInfo structure
		purchaseInfo := types.PurchaseInfo{
			FirstName:    fname,
			LastName:     lname,
			Email:        email,
			Car:          product.Name,
			Quantity:     quantity,
			Price:        product.Price,
			Total:        totalWithTax,
			RoundUpTotal: grandtotal,
			Message:      welcomeMessage,
		}

		// Add an order to the database
		orderID, err := db.AddOrder(conn, types.Order{
			CustomerFirstName: fname,
			CustomerLastName:  lname,
			ProductName:       product.Name,
			Quantity:          quantity,
			Price:             grandtotal,
			Tax:               tax,
			Donation:          donation,
			Timestamp:         timestamp,
		})
		if err != nil {
			fmt.Printf("Error adding order: %v\n", err)
			return ctx.String(http.StatusInternalServerError, "Error adding order")
		}

		fmt.Printf("Order added with ID: %v\n", orderID)

		return Render(ctx, http.StatusOK, templates.Base(templates.PurchaseConfirmation(purchaseInfo)))
	})

	e.GET("/order_entry", func(ctx echo.Context) error {
		products, err := db.GetAllProducts(conn)
		if err != nil {
			return ctx.String(http.StatusInternalServerError, fmt.Sprintf("Error loading products: %v", err))
		}

		return Render(ctx, http.StatusOK, templates.OrderEntry(products))
	})

	e.GET("/get_customers", func(ctx echo.Context) error {
		query := ctx.QueryParam("query")
		field := ctx.QueryParam("field")

		//log.Printf("Received query: %s, field: %s", query, field)

		if query == "" || (field != "first" && field != "last") {
			log.Println("Invalid query parameters")
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid query parameters"})
		}

		var customers []types.Customer

		var rows *sql.Rows
		var err error

		if field == "first" {
			// Get customers by fname
			rows, err = conn.Query("SELECT first_name, last_name, email FROM customer WHERE first_name LIKE ?", query+"%")
		} else {
			// Get customers by lname
			rows, err = conn.Query("SELECT first_name, last_name, email FROM customer WHERE last_name LIKE ?", query+"%")
		}

		if err != nil {
			log.Printf("Database query error: %v", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Database query error"})
		}
		defer rows.Close()

		for rows.Next() {
			var customer types.Customer
			if err := rows.Scan(&customer.FirstName, &customer.LastName, &customer.Email); err != nil {
				log.Printf("Row scan error: %v", err)
				return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Row scan error"})
			}
			customers = append(customers, customer)
		}

		if rows.Err() != nil {
			log.Printf("Rows error after iteration: %v", rows.Err())
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Error iterating rows"})
		}

		//log.Printf("Returning customers: %+v", customers)

		return ctx.JSON(http.StatusOK, map[string]interface{}{"customers": customers})
	})

	e.POST("/add_order", func(ctx echo.Context) error {
		customerName := ctx.FormValue("fname") + " " + ctx.FormValue("lname")
		productName := ctx.FormValue("car")
		quantity, _ := strconv.Atoi(ctx.FormValue("quantity"))

		// get stock from db
		var availableQuantity int
		var pricePerUnit float64
		err := conn.QueryRow("SELECT in_stock, price FROM product WHERE product_name = ?", productName).Scan(&availableQuantity, &pricePerUnit)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not retrieve product details"})
		}

		// make sure they dont buy more than available
		if quantity > availableQuantity {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Requested quantity exceeds available stock",
			})
		}

		total := float64(quantity) * pricePerUnit

		// update db
		_, err = conn.Exec("UPDATE product SET in_stock = in_stock - ? WHERE product_name = ?", quantity, productName)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not complete the order"})
		}

		// confirm msg string
		confirmationMessage := fmt.Sprintf("Order submitted for %s: %d x %s for a total of $%.2f", customerName, quantity, productName, total)

		return ctx.JSON(http.StatusOK, map[string]string{"message": confirmationMessage})
	})

	e.GET("/admin", func(ctx echo.Context) error {
		// Get all customers
		customers, err := db.GetAllCustomers(conn)
		if err != nil {
			e.Logger.Errorf("Error fetching customers: %v", err)
			return ctx.String(http.StatusInternalServerError, "Error loading customers")
		}

		// Get all orders
		orders, err := db.GetAllOrders(conn)
		if err != nil {
			e.Logger.Errorf("Error fetching orders: %v", err)
			return ctx.String(http.StatusInternalServerError, "Error loading orders")
		}

		// Get all products
		products, err := db.GetAllProducts(conn)
		if err != nil {
			e.Logger.Errorf("Error fetching products: %v", err)
			return ctx.String(http.StatusInternalServerError, "Error loading products")
		}

		// Render the admin page
		return Render(ctx, http.StatusOK, templates.Base(templates.AdminPage(customers, orders, products)))
	})

	e.GET("/products", func(ctx echo.Context) error {
		// Get all products
		products, err := db.GetAllProducts(conn)

		for _, product := range products {
			fmt.Printf("Product: %s, Inactive: %d\n", product.Name, product.Inactive)
		}

		if err != nil {
			e.Logger.Errorf("Error fetching products: %v", err)
			return ctx.String(http.StatusInternalServerError, "Error loading products")
		}

		return Render(ctx, http.StatusOK, templates.Base(templates.Products(products)))
	})

	e.POST("/add_car", func(ctx echo.Context) error {
		var product types.Product

		product.Name = ctx.FormValue("name")
		product.Image = ctx.FormValue("image")

		price, err := strconv.ParseFloat(ctx.FormValue("price"), 64)
		if err != nil {
			return ctx.String(http.StatusBadRequest, "Invalid price value")
		}
		product.Price = price

		quantity, err := strconv.Atoi(ctx.FormValue("quantity"))
		if err != nil {
			return ctx.String(http.StatusBadRequest, "Invalid quantity value")
		}
		product.QuantityInStock = quantity

		inactive, _ := strconv.Atoi(ctx.FormValue("inactive"))
		product.Inactive = inactive

		// Add car to the db
		if err := db.AddCar(conn, product); err != nil {
			return ctx.String(http.StatusInternalServerError, fmt.Sprintf("Error adding car: %v", err))
		}

		return ctx.JSON(http.StatusOK, map[string]string{"message": "Car added successfully"})
	})

	e.PUT("/update_car/:id", func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return ctx.String(http.StatusBadRequest, "Invalid car ID")
		}
	
		var product types.Product
		product.ID = int64(id)
	
		// Get form values
		product.Name = ctx.FormValue("name")
		product.Image = ctx.FormValue("image")
	
		// Get price
		priceStr := ctx.FormValue("price")
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return ctx.String(http.StatusBadRequest, "Invalid price value")
		}
		product.Price = price
	
		// Get quantity
		quantityStr := ctx.FormValue("quantity")
		quantity, err := strconv.Atoi(quantityStr)
		if err != nil {
			return ctx.String(http.StatusBadRequest, "Invalid quantity value")
		}
		product.QuantityInStock = quantity
	
		// Retrieve the inactive field, defaulting to 0 if missing
		inactiveStr := ctx.FormValue("inactive")
		inactive := 0
		if inactiveStr == "1" {
			inactive = 1
		}
		product.Inactive = inactive
	
		fmt.Printf("Updating product: %+v\n", product) // Debugging line
	
		// Update the car in the database
		if err := db.UpdateCar(conn, product); err != nil {
			return ctx.String(http.StatusInternalServerError, fmt.Sprintf("Error updating car: %v", err))
		}
	
		return ctx.JSON(http.StatusOK, map[string]string{"message": "Car updated successfully"})
	})
	

	e.DELETE("/delete_car/:id", func(ctx echo.Context) error {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			return ctx.String(http.StatusBadRequest, "Invalid car ID")
		}

		if err := db.DeleteCar(conn, int64(id)); err != nil {
			return ctx.String(http.StatusInternalServerError, fmt.Sprintf("Error deleting car: %v", err))
		}
		return ctx.JSON(http.StatusOK, map[string]string{"message": "Car deleted successfully"})
	})

	e.GET("/get_all_products", func(ctx echo.Context) error {
		products, err := db.GetAllProducts(conn)
		if err != nil {
			e.Logger.Errorf("Error fetching products: %v", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Error loading products"})
		}
		return ctx.JSON(http.StatusOK, map[string]interface{}{"products": products})
	})

	e.Logger.Fatal(e.Start(":8000"))
}

// INFO: This is a simplified render method that replaces `echo`'s with a custom
// one. This should simplify rendering out of an echo route.
func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}
