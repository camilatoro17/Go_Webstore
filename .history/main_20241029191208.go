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

	// Get a database handle.
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
		// TODO: Grab the form details from ctx.FormValue("...")

		fname := ctx.FormValue("fname")
		lname := ctx.FormValue("lname")
		email := ctx.FormValue("email")
		carIDStr := ctx.FormValue("car")
		quantitystr := ctx.FormValue("quantity")
		roundup := ctx.FormValue("donate")

		carID, err := strconv.ParseInt(carIDStr, 10, 64)
    if err != nil {
        return ctx.String(http.StatusBadRequest, "Invalid car ID")
    }

		// Fetch product details from the database based on car name
		products, err := db.ProductByName(conn, car)
		if err != nil || len(products) == 0 {
			return ctx.String(http.StatusBadRequest, "Car not found")
		}
		product := products[0]

		// Convert quantity from str to int
		quantity, err := strconv.Atoi(quantitystr)
		if err != nil {
			return ctx.String(http.StatusBadRequest, "Invalid quantity")
		}

		// Update the product quantity after the sale
		err = db.UpdateProductQuantity(conn, car, quantity)
		if err != nil {
			log.Printf("Error updating product quantity: %v", err)
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

		// Create the purchaseInfo structure
		purchaseInfo := types.PurchaseInfo{
			FirstName:    fname,
			LastName:     lname,
			Email:        email,
			Car:          car,
			Quantity:     quantity,
			Price:        product.Price,
			Total:        totalWithTax,
			RoundUpTotal: grandtotal,
		}

		// Add an order to the database
		orderID, err := db.AddOrder(conn, types.Order{
			CustomerFirstName: fname,
			CustomerLastName:  lname,
			ProductName:       car,
			Quantity:          quantity,
			Price:             grandtotal,
			Tax:               tax,
			Donation:          donation,
		})
		if err != nil {
			log.Printf("Error adding order: %v", err)
			return ctx.String(http.StatusInternalServerError, "Error adding order")
		}

		fmt.Printf("Order added with ID: %v\n", orderID)

		return Render(ctx, http.StatusOK, templates.Base(templates.PurchaseConfirmation(purchaseInfo)))
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
