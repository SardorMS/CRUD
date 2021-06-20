package customers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/SardorMS/CRUD/pkg/types"
)

var (
	ErrInternal        = errors.New("internal error") //return when an internal error occurred
	ErrNoSuchUser      = errors.New("no such user")
	ErrPhoneUsed       = errors.New("phone already registred")
	ErrInvalidPassword = errors.New("invalid password")
	ErrTokenNotFound   = errors.New("token not found") //retrun when customer not found
	ErrTokenExpired    = errors.New("token expired")
	ErrNotFound        = errors.New("not found")
)

//Service - describes customer service.
type Service struct {
	pool *pgxpool.Pool
}

//newService - create a service.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

// IDByToken - performs the users authentication procedure,
// if successfull returns its id.
// Returns ErrNoSuchUser if user is not found.
// Returns ErrInternal if another error occurs.
func (s *Service) IDByToken(ctx context.Context, token string) (int64, error) {

	var id int64
	sql := `SELECT customer_id FROM customers_tokens WHERE token = $1;`
	err := s.pool.QueryRow(ctx, sql, token).Scan(&id)

	if err == pgx.ErrNoRows {
		return 0, nil
	}

	if err != nil {
		return 0, ErrInternal
	}
	return id, nil
}

//Tokern - generates a token for user.
// Returns ErrNoSuchUser if user is not found.
// Returns ErrInvalidPassword if password incorrect.
// Returns ErrInternal if another error occurs.
func (s *Service) Token(ctx context.Context, phone string, password string,
) (token string, err error) {

	var id int64
	var hash string

	sql1 := `SELECT id, password FROM customers WHERE phone = $1;`
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
	sql2 := `INSERT INTO customers_tokens(token, customer_id) VALUES($1, $2);`
	_, err = s.pool.Exec(ctx, sql2, token, id)
	if err != nil {
		return "", ErrInternal
	}
	return token, nil
}

// Register - ...
func (s *Service) Register(ctx context.Context, registration *types.Registration) (*types.Customer, error) {

	var err error
	item := &types.Customer{}

	hash, err := bcrypt.GenerateFromPassword([]byte(registration.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return nil, ErrPhoneUsed
	}

	sql1 := `INSERT INTO customers (name, phone, password) 
		VALUES ($1, $2, $3) ON CONFLICT (phone) DO NOTHING 
		RETURNING id, name, phone, password, active, created;`
	err = s.pool.QueryRow(ctx, sql1, registration.Name, registration.Phone, hash).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
		&item.Active,
		&item.Created,
	)

	if err == pgx.ErrNoRows {
		log.Println(err)
		return nil, ErrNoSuchUser
	}
	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}
	return item, nil
}

// Products - ...
func (s *Service) Products(ctx context.Context) ([]*types.Product, error) {

	items := make([]*types.Product, 0)
	sql := `SELECT id, name, price, qty FROM products WHERE active ORDER BY id LIMIT 500;`
	rows, err := s.pool.Query(ctx, sql)

	if errors.Is(err, pgx.ErrNoRows) {
		return items, nil
	}

	if err != nil {
		return nil, ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		item := &types.Product{}
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

//Purchases -
func (s *Service) Purchases(ctx context.Context, id int64) ([]*types.Sales, error) {

	items := make([]*types.Sales, 0)

	sql := `SELECT sp.id, sp.name, sp.price, sp.qty, sp.created 
			FROM sale_positions sp
			JOIN sales s ON s.id = sp.sale_id
			WHERE s.customers_id = $1;
			`
	rows, err := s.pool.Query(ctx, sql, id)

	if errors.Is(err, pgx.ErrNoRows) {
		return items, nil
	}

	if err != nil {
		return nil, ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		item := &types.Sales{}
		err = rows.Scan(
			&item.ID,
			&item.Name,
			&item.Price,
			&item.Qty,
			&item.Created)

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

//All - returns all existing customers.
func (s *Service) All(ctx context.Context) ([]*types.Customer, error) {

	items := make([]*types.Customer, 0)
	sql := `SELECT id, name, phone, active, created FROM customers;`
	rows, err := s.pool.Query(ctx, sql)
	if err != nil {
		log.Println(err)
		return nil, ErrNotFound
	}
	defer rows.Close()

	for rows.Next() {
		item := &types.Customer{}
		err = rows.Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created,
		)

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

	log.Println(items)
	return items, nil
}

//AllActive - returns all existing customers.
func (s *Service) AllActive(ctx context.Context) ([]*types.Customer, error) {

	items := make([]*types.Customer, 0)
	sql := `SELECT id, name, phone, active, created FROM customers WHERE active = TRUE;`
	rows, err := s.pool.Query(ctx, sql)
	if err != nil {
		log.Println(err)
		return nil, ErrNotFound
	}
	defer rows.Close()

	for rows.Next() {
		item := &types.Customer{}
		err = rows.Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created,
		)
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

	log.Println(items)
	return items, nil
}

//Save - creates, saves and updates customer lists.
func (s *Service) Save(ctx context.Context, items *types.Customer) (*types.Customer, error) {

	item := &types.Customer{}

	if items.ID == 0 {
		sql1 := `INSERT INTO customers (name, phone) VALUES ($1, $2, $3) ON CONFLICT (phone) DO UPDATE SET name= excluded.name
		RETURNING id, name, phone, active, created;`
		err := s.pool.QueryRow(ctx, sql1, items.Name, items.Phone).Scan(
			&item.ID,
			&item.Name,
			&item.Phone,
			&item.Active,
			&item.Created,
		)
		if err != nil {
			log.Println(err)
			return nil, ErrInternal
		}
		return item, nil
	}

	sql2 := `UPDATE customers SET name = $2, phone = $3 WHERE id = $1 
	RETURNING  id, name, phone, active, created;`
	err := s.pool.QueryRow(ctx, sql2, items.ID, items.Name, items.Phone).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
		&item.Active,
		&item.Created,
	)

	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}
	return item, nil
}

//Delete - deletes customer lists.
func (s *Service) Delete(ctx context.Context, id int64) (*types.Customer, error) {
	item := &types.Customer{}

	sql := `DELETE FROM customers WHERE id = $1 RETURNING id, name, phone, active, created;`
	err := s.pool.QueryRow(ctx, sql, id).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
		&item.Active,
		&item.Created,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Println("No Rows")
		return nil, ErrNotFound
	}

	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}
	return item, nil
}

//Block - block customers list(sets false status).
func (s *Service) Block(ctx context.Context, id int64, active bool) (*types.Customer, error) {
	item := &types.Customer{}

	sql := `UPDATE customers SET active = $2 WHERE id = $1 RETURNING id, name, phone, active, created;`
	err := s.pool.QueryRow(ctx, sql, id, active).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
		&item.Active,
		&item.Created,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Println("No Rows")
		return nil, ErrNotFound
	}

	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}
	return item, nil
}

//Unblock - unblock customers list(sets true status).
func (s *Service) Unblock(ctx context.Context, id int64, active bool) (*types.Customer, error) {
	item := &types.Customer{}

	sql := `UPDATE customers SET active = $2 WHERE id = $1 RETURNING id, name, phone, active, created;`
	err := s.pool.QueryRow(ctx, sql, id, active).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
		&item.Active,
		&item.Created,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Println("No Rows")
		return nil, ErrNotFound
	}

	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}
	return item, nil
}

//ByID - returns customer by ID.
func (s *Service) ByID(ctx context.Context, id int64) (*types.Customer, error) {
	item := &types.Customer{}

	sql := `SELECT id, name, phone, active, created FROM customers WHERE id = $1;`
	err := s.pool.QueryRow(ctx, sql, id).Scan(
		&item.ID,
		&item.Name,
		&item.Phone,
		&item.Active,
		&item.Created,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		log.Println("No Rows")
		return nil, ErrNotFound
	}

	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}
	return item, nil
}