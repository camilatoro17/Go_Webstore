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

	// Capture connection properties.
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

		// Get products from the database
		products, err := db.GetAllProducts(conn)
		if err != nil {
			e.Logger.Errorf("Error fetching products: %v", err)
			return ctx.String(http.StatusInternalServerError, "Error loading products")
		}

		// Render the store page with the products
		return Render(ctx, http.StatusOK, templates.Base(templates.Store(products)))
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
		carIDStr := ctx.FormValue("car") // Get the car ID
		quantityStr := ctx.FormValue("quantity")
		roundup := ctx.FormValue("donate")
		timestampStr := ctx.FormValue("timestamp")

		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			fmt.Printf("Invalid timestamp: %v\n", err)
			return ctx.String(http.StatusBadRequest, "Invalid timestamp")
		}

		// Check if the order with the same timestamp already exists
		exists, err := db.OrderExists(conn, timestamp)
		if err != nil {
			log.Printf("Error checking order existence: %v", err)
			return ctx.String(http.StatusInternalServerError, "Error checking order existence")
		}

		if exists {
			fmt.Println("Duplicate order detected. Skipping insertion.")
			return ctx.String(http.StatusOK, "This order has already been processed.")
		}


		// Convert car ID from string to int64
		carID, err := strconv.ParseInt(carIDStr, 10, 64)
		if err != nil {
			fmt.Printf("Invalid car ID: %v\n", err)
			return ctx.String(http.StatusBadRequest, "Invalid car ID")
		}

		// Fetch product details from the database using the car ID
		product, err := db.ProductByID(conn, carID)
		if err != nil {
			fmt.Printf("ProductByID error: %v\n", err)
			return ctx.String(http.StatusNotFound, "Car not found")
		}

		// Convert quantity from string to int
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

	e.GET("/admin", func(ctx echo.Context) error {
		// Get all customers
		customers, err := db.GetAllCustomers(conn)
		if err != nil {
			e.Logger.Errorf("Error fetching customers: %v", err)
			return ctx.String(http.StatusInternalServerError, "Error loading customers")
		}
	
		// Get all orders (now using GetAllOrders)
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
