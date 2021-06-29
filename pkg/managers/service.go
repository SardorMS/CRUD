package managers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"strconv"

	"github.com/SardorMS/CRUD/pkg/types"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInternal        = errors.New("internal error")          //return when an internal error occurred.
	ErrNoSuchUser      = errors.New("no such user")            //return no such user.
	ErrPhoneUsed       = errors.New("phone already registred") //retunr error if phone is already registred.
	ErrInvalidPassword = errors.New("invalid password")        //return invalid password.
	ErrTokenNotFound   = errors.New("token not found")         //retrun when token not found.
	ErrTokenExpired    = errors.New("token expired")           //return when token expired
	ErrNotFound        = errors.New("not found")               // return not found
)

//Service - describes managers service.
type Service struct {
	pool *pgxpool.Pool
}

//newService - create a service.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// IDByToken - performs the users authentication procedure,
// if successfull returns its id.
func (s *Service) IDByToken(ctx context.Context, token string) (int64, error) {

	var id int64
	sql := `SELECT manager_id FROM managers_tokens WHERE token = $1;`
	err := s.pool.QueryRow(ctx, sql, token).Scan(&id)

	if err != nil {
		log.Println(err)
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, nil
	}
	return id, nil
}

// IsAdmin - managers role check.
func (s *Service) IsAdmin(ctx context.Context, id int64) (isAdmin bool) {
	sql := `SELECT is_admin FROM managers WHERE id = $1`
	err := s.pool.QueryRow(ctx, sql, id).Scan(&isAdmin)
	if err != nil {
		log.Println(err)
		return false
	}
	return
}

// Register - managers register procedure.
func (s *Service) Register(ctx context.Context, item *types.Managers) (string, error) {

	var token string
	var id int64

	sql1 := `INSERT INTO managers (name, phone, is_admin) 
		VALUES ($1, $2, $3) ON CONFLICT (phone) DO NOTHING RETURNING id;`
	err := s.pool.QueryRow(ctx, sql1, item.Name, item.Phone, item.IsAdmin).Scan(&id)
	if err != nil {
		log.Print(err)
		return "", ErrInternal
	}

	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		return "", ErrInternal
	}
	token = hex.EncodeToString(buffer)

	sql2 := `INSERT INTO managers_tokens (token, manager_id) VALUES($1, $2);`
	_, err = s.pool.Exec(ctx, sql2, token, id)
	if err != nil {
		log.Println(err)
		return "", ErrInternal
	}
	return token, nil
}

// Tokern - generates a token for managers.
func (s *Service) Token(ctx context.Context, phone string, password string) (token string, err error) {

	var id int64
	var hash string

	sql1 := `SELECT id, password FROM managers WHERE phone = $1;`
	err = s.pool.QueryRow(ctx, sql1, phone).Scan(&id, &hash)

	if err == pgx.ErrNoRows {
		return "", ErrNoSuchUser
	}

	if err != nil {
		return "", ErrInternal
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return "", ErrInvalidPassword
	}

	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil {
		return "", ErrInternal
	}

	token = hex.EncodeToString(buffer)
	sql2 := `INSERT INTO managers_tokens (token, manager_id) VALUES($1, $2);`
	_, err = s.pool.Exec(ctx, sql2, token, id)
	if err != nil {
		return "", ErrInternal
	}
	return token, nil
}

// GetSales - get the sales information.
func (s *Service) GetSales(ctx context.Context, id int64) (sum int, err error) {

	sql := `SELECT COALESCE (SUM (sp.price * sp.qty), 0) total
			FROM managers m
			LEFT JOIN sales s ON s.manager_id = $1
			LEFT JOIN sale_positions sp ON sp.sale_id = s.id
			GROUP BY m.id
			LIMIT 1;`
	err = s.pool.QueryRow(ctx, sql, id).Scan(&sum)

	if err != nil {
		log.Println(err)
		return 0, ErrInternal
	}
	return sum, nil
}

// MakeSales - makes a sale.
func (s *Service) MakeSales(ctx context.Context, sale *types.Sale) (*types.Sale, error) {

	positionSQL := "INSERT INTO sale_positions (sale_id, product_id, qty, price) VALUES "

	sql := `INSERT INTO sales (manager_id, customer_id) VALUES ($1, $2) RETURNING id, created;`
	err := s.pool.QueryRow(ctx, sql, sale.ManagerID, sale.CustomerID).Scan(&sale.ID, &sale.Created)

	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}

	for _, position := range sale.Positions {
		if !s.MakeSalePosition(ctx, position) {
			log.Println("Invalid position")
			return nil, ErrInternal
		}
		positionSQL += "(" + strconv.FormatInt(sale.ID, 10) + "," +
			strconv.FormatInt(position.ProductID, 10) + "," +
			strconv.Itoa(position.Price) + "," + strconv.Itoa(position.Qty) + "),"
	}

	positionSQL = positionSQL[0 : len(positionSQL)-1]
	log.Println(positionSQL)

	_, err = s.pool.Exec(ctx, positionSQL)
	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}

	return sale, nil
}

// MakeSalePosition - updates a products information.
func (s *Service) MakeSalePosition(ctx context.Context, position *types.SalePosition) bool {
	active := false
	qty := 0

	sql1 := `SELECT qty, active FROM products WHERE id = $1;`
	if err := s.pool.QueryRow(ctx, sql1, position.ProductID).Scan(&qty, &active); err != nil {
		return false
	}

	if qty < position.Qty || !active {
		return false
	}

	sql2 := `UPDATE products SET qty = $1 WHERE id = $2`
	if _, err := s.pool.Exec(ctx, sql2, qty-position.Qty, position.ProductID); err != nil {
		log.Println(err)
		return false
	}
	return true
}

// Products - shows information about products to customers.
func (s *Service) Products(ctx context.Context) ([]*types.Products, error) {

	items := make([]*types.Products, 0)
	sql := `SELECT id, name, price, qty FROM products WHERE active = true ORDER BY id LIMIT 500;`
	rows, err := s.pool.Query(ctx, sql)

	if errors.Is(err, pgx.ErrNoRows) {
		return items, nil
	}

	if err != nil {
		return nil, ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		item := &types.Products{}
		err = rows.Scan(&item.ID, &item.Name, &item.Price, &item.Qty)

		if err != nil {
			log.Println(err)
			return nil, err
		}
		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return items, nil
}

// ChangeProduct(Save) - change or save an information about products.
func (s *Service) ChangeProduct(ctx context.Context, product *types.Products) (*types.Products, error) {

	var err error

	if product.ID == 0 {
		sql1 := `INSERT INTO products (name, qty, price) VALUES ($1, $2, $3)
				 RETURNING id, name, qty, price, active, created;`
		err = s.pool.QueryRow(ctx, sql1, product.Name, product.Qty, product.Price).Scan(
			&product.ID,
			&product.Name,
			&product.Qty,
			&product.Price,
			&product.Active,
			&product.Created)

	} else {
		sql2 := `UPDATE products SET name = $1, qty = $2, price = $3 WHERE id = $4 
				 RETURNING  id, name, qty, price, active, created;`
		err = s.pool.QueryRow(ctx, sql2, product.Name, product.Qty, product.Price, product.ID).Scan(
			&product.ID,
			&product.Name,
			&product.Qty,
			&product.Price,
			&product.Active,
			&product.Created)
	}

	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}
	return product, nil

}

// RemoveProductByID - remove information about products.
func (s *Service) RemoveProductByID(ctx context.Context, id int64) (*types.Products, error) {
	item := &types.Products{}

	sql := `DELETE FROM products WHERE id = $1 RETURNING id, name, price, qty, active, created;`
	err := s.pool.QueryRow(ctx, sql, id).Scan(
		&item.ID,
		&item.Name,
		&item.Price,
		&item.Qty,
		&item.Active,
		&item.Created)

	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}
	return item, nil
}

// GetCustomer -  shows information about customers.
func (s *Service) GetCustomer(ctx context.Context) ([]*types.Customers, error) {
	items := make([]*types.Customers, 0)
	sql := `SELECT id, name, phone, active, created FROM customers 
			WHERE active = true ORDER BY id LIMIT 500;`

	rows, err := s.pool.Query(ctx, sql)

	if err != nil {
		if err == pgx.ErrNoRows {
			return items, nil
		}

		return nil, ErrNotFound
	}
	defer rows.Close()

	for rows.Next() {
		item := &types.Customers{}
		err = rows.Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created)

		if err != nil {
			log.Println(err)
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// ChangeCustomer - change information about customers.
func (s *Service) ChangeCustomer(ctx context.Context, customer *types.Customers) (*types.Customers, error) {

	sql := `UPDATE customers SET name = $1, phone = $2, active = $3 WHERE id = $4 
			RETURNING name, phone, active;`
	err := s.pool.QueryRow(ctx, sql, customer.Name, customer.Phone, customer.Active, customer.ID).Scan(
		&customer.Name,
		&customer.Phone,
		&customer.Active)

	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}

	return customer, nil
}

// RemoveCustomerByID - remove information about the ucustomers ID.
func (s *Service) RemoveCustomerByID(ctx context.Context, id int64) (*types.Customers, error) {
	item := &types.Customers{}

	sql := `DELETE FROM customers WHERE id = $1 RETURNING id, name, phone, active, created;`
	err := s.pool.QueryRow(ctx, sql, id).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
		&item.Active,
		&item.Created)

	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}
	return item, nil
}
