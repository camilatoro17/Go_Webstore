package types

// TODO: If you choose to use a struct rather than individual parameters to your view, you might flesh this one out:
type PurchaseInfo struct {
	FirstName    string
	LastName     string
	Email        string
	Car          string
	Quantity     int
	Price        float64
	Total        float64
	RoundUpTotal float64
}

type CustomerResults struct {
	Customers []Customer
	Customer2 Customer
	Customer3 Customer
	Customer1 Customer
}

type Customer struct {
	ID        id
	FirstName string
	LastName  string
	Email     string
}
