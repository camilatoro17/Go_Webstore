package main

import (
	"net/http"
	"strconv"

	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/a-h/templ"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	etag "github.com/pablor21/echo-etag/v4"

	"go-store/templates"
	"go-store/types"
)

var products = map[string]float64{
	"2024 G80 M3":   78000,
	"2024 S63 AMG":  183000,
	"2024 Audi RS7": 128000,
}

type Product struct {
    ID       int64
    Name     string
    Image    string
    Price    float64
    InStock  int
}

type Customer struct {
    ID       int64
    FirstName string
    LastName  string
    Email     string
}

type Order struct {
    ID         int64
    ProductID  int
    CustomerID int
    Quantity   int
    Price      float64
}

var db *sql.DB

func main() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "ctorosuarez",
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	// TODO: Fill in your products here with name -> price as the key -> value pair.
	// products := map[string]float64{
	// }
	e := echo.New()
	e.Use(etag.Etag())

	// INFO: If you wanted to load a CSS file, you'd do something like this:
	// `<link rel="stylesheet" href="assets/styles/styles.css">`
	e.Static("assets", "./assets")

	// TODO: Render your base store page here
	e.GET("/store", func(ctx echo.Context) error {
		return Render(ctx, http.StatusOK, templates.Base(templates.Store(products)))
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
		products, err := productByName("2024 G80 M3")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Products found: %v\n", products)

		//find customers by last name
		customers, err := customerByLastName("Mouse")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Customers found: %v\n", customers)

		//find orders by customer id=1
		orders, err := ordersByCustomer(1)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Orders found: %v\n", orders)

		//find a product by its id=1
		p, err := productByID(1)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Product found: %v\n", p)

		//add an order
		orderID, err := addOrder(Order{
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

func addOrder(ord Order) (int64, error) {
    result, err := db.Exec("INSERT INTO orders (product_id, customer_id, quantity, price) VALUES (?, ?, ?, ?)",
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


func productByName(name string) ([]Product, error) {
    var products []Product
    rows, err := db.Query("SELECT * FROM product WHERE product_name = ?", name)
    if err != nil {
        return nil, fmt.Errorf("productByName %q: %v", name, err)
    }
    defer rows.Close()

    for rows.Next() {
        var p Product
        if err := rows.Scan(&p.ID, &p.Name, &p.Image, &p.Price, &p.InStock); err != nil {
            return nil, fmt.Errorf("productByName %q: %v", name, err)
        }
        products = append(products, p)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return products, nil
}

func productByID(id int64) (Product, error) {
    var p Product
    row := db.QueryRow("SELECT * FROM product WHERE id = ?", id)
    
    err := row.Scan(&p.ID, &p.Name, &p.Image, &p.Price, &p.InStock)
    if err != nil {
        if err == sql.ErrNoRows {
            return p, fmt.Errorf("productByID %d: no such product", id)
        }
        return p, fmt.Errorf("productByID %d: %v", id, err)
    }

    return p, nil
}

func customerByLastName(lastName string) ([]Customer, error) {
    var customers []Customer
    rows, err := db.Query("SELECT * FROM customer WHERE last_name = ?", lastName)
    if err != nil {
        return nil, fmt.Errorf("customerByLastName %q: %v", lastName, err)
    }
    defer rows.Close()

    for rows.Next() {
        var c Customer
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

func ordersByCustomer(customerID int) ([]Order, error) {
    var orders []Order
    rows, err := db.Query("SELECT * FROM orders WHERE customer_id = ?", customerID)
    if err != nil {
        return nil, fmt.Errorf("ordersByCustomer %d: %v", customerID, err)
    }
    defer rows.Close()

    for rows.Next() {
        var o Order
        if err := rows.Scan(&o.ID, &o.ProductID, &o.CustomerID, &o.Quantity, &o.Price); err != nil {
            return nil, fmt.Errorf("ordersByCustomer %d: %v", customerID, err)
        }
        orders = append(orders, o)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return orders, nil
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
