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

var products = map[string]float64{
	"2024 G80 M3":   78000,
	"2024 S63 AMG":  183000,
	"2024 Audi RS7": 128000,
}

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
		return Render(ctx, http.StatusOK, templates.Base(templates.Store(products)))
	})

	e.GET("/schema", func(ctx echo.Context) error {
		products, err := db.GetAllProducts(conn)
		e.Logger.Errorf("%+v", err)
		return Render(ctx, http.StatusOK, templates.Queries(products))
	})

	// TODO: Handle the form submission and return the purchase confirmation view
	e.POST("/purchase", func(ctx echo.Context) error {
		// TODO: Grab the form details from ctx.FormValue("...")

		fname := ctx.FormValue("fname")
		lname := ctx.FormValue("lname")
		email := ctx.FormValue("email")
		car := ctx.FormValue("car")
		quantitystr := ctx.FormValue("quantity")
		roundup := ctx.FormValue("donate")

		price, exists := products[car]
		if !exists {
			return ctx.String(http.StatusBadRequest, "Car not found")
		}

		// Convert quantity from str to int
		quantity, err := strconv.Atoi(quantitystr)
		if err != nil {
			return ctx.String(http.StatusBadRequest, "Invalid quantity")
		}

		// Calculate subtotal (price * quantity)
		subtotal := price * float64(quantity)

		// Calculate tax (2.9%)
		const taxRate = 0.029
		tax := subtotal * taxRate

		// Calculate total including tax
		totalWithTax := subtotal + tax

		// Check if user opted to round up
		var grandtotal float64
		if roundup == "yes" {
			grandtotal = totalWithTax + 1.00 // Round up
		} else {
			grandtotal = totalWithTax
		}

		// TODO: Maybe use this structure to pass the data to your purchase confirmation page
		// ...
		purchaseInfo := types.PurchaseInfo{
			FirstName:    fname,
			LastName:     lname,
			Email:        email,
			Car:          car,
			Quantity:     quantity,
			Price:        price,
			Total:        totalWithTax,
			RoundUpTotal: grandtotal,
		}

		//find products by the name
		products, err := db.ProductByName(conn, "2024 G80 M3")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Products found: %v\n", products)

		//find customers by last name
		customers, err := db.CustomerByLastName(conn, "Mouse")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Customers found: %v\n", customers)

		//find orders by customer id=1
		orders, err := db.OrdersByCustomer(conn, 1)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Orders found: %v\n", orders)

		//find a product by its id=1
		p, err := db.ProductByID(conn, 1)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Product found: %v\n", p)

		//add an order
		orderID, err := db.AddOrder(conn, types.Order{
			ProductID:  1,
			CustomerID: 1,
			Quantity:   2,
			Price:      183000,
		})
		if err != nil {
			log.Fatal(err)
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
