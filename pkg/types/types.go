package types

import "time"

// Customer - represents information about customer.
type Customer struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
}

// Registration -
type Registration struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// Auth - ...
type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Token - ...
type Token struct {
	Token string `json:"token"`
}

// Product - ...
type Product struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Qty   int    `json:"qty"`
}

//Purchases (sales) - ...
type Sales struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Price   int       `json:"price"`
	Qty     int       `json:"qty"`
	Created time.Time `json:"created"`
}


//       
//-------------------------------------------------------------------------//
//


// ManagerRegister - ...
type ManagerRegister struct {
	ID    int64    `json:"id"`
	Name  string   `json:"name"`
	Phone string   `json:"phone"`
	Roles []string `json:"roles"`
}

//Managers - represents information about customers.
type Managers struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	Password   string    `json:"password"`
	Salary     int64     `json:"salary"`
	Plan       int64     `json:"plan"`
	BossID     int64     `json:"boss_id"`
	Department string    `json:"department"`
	IsAdmin    bool      `json:"is_admin"`
	Created    time.Time `json:"created"`
}

// Sale - ...
type Sale struct {
	ID         int64           `json:"id"`
	ManagerID  int64           `json:"manager_id"`
	CustomerID int64           `json:"customer_id"`
	Created    time.Time       `json:"created"`
	Positions  []*SalePosition `json:"positions"`
}

// SalePositions - ...
type SalePosition struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	SaleID    int64     `json:"sale_id"`
	Price     int       `json:"price"`
	Qty       int       `json:"qty"`
	Created   time.Time `json:"created"`
}

// - Products - ...
type Products struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Price   int       `json:"price"`
	Qty     int       `json:"qty"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
}

// Customeres - ...
type Customers struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Phone   string    `json:"phone"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
}