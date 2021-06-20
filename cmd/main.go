package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/dig"

	"github.com/SardorMS/CRUD/cmd/app"
	"github.com/SardorMS/CRUD/pkg/customers"
	"github.com/SardorMS/CRUD/pkg/security"
)

func main() {

	host := "127.0.0.1"
	port := "9999"
	//user:login@host:port/db
	dsn := "postgres://sardor:123@192.168.99.100:5432/db"

	if err := execute(host, port, dsn); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(host string, port string, dsn string) (err error) {

	deps := []interface{}{
		app.NewServer,
		mux.NewRouter,
		func() (*pgxpool.Pool, error) {
			ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
			return pgxpool.Connect(ctx, dsn)
		},
		customers.NewService,
		security.NewService,
		//managers.NewService,
		//products.NewService,
		//sales.NewService,
		func(server *app.Server) *http.Server {
			return &http.Server{
				Addr:    net.JoinHostPort(host, port),
				Handler: server,
			}
		},
	}

	container := dig.New()
	for _, dep := range deps {
		err = container.Provide(dep)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	err = container.Invoke(func(server *app.Server) {
		server.Init()
	})
	if err != nil {
		log.Println(err)
		return err
	}

	return container.Invoke(func(server *http.Server) error {
		return server.ListenAndServe()
	})

}

/*
//Create
	ctx := context.Background()
	_, err = db.ExecContext(ctx, `
	  CREATE TABLE IF NOT EXISTS customers
	(
		id      BIGSERIAL PRIMARY KEY,
    	name    TEXT      NOT NULL,
    	phone   TEXT      NOT NULL CHECK UNIQUE,
    	active  BOOLEAN   NOT NULL DEFAULT TRUE,
    	created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)
	`)
	if err != nil {
		log.Println(err)
		//os.Exit(1) делать нельзя, иначе defer не выполнится
		return
	}

	//Insert
	name := "Vasya"
	phone := "+992000000001"
	resutl, err := db.ExecContext(ctx, `
	  INSERT INTO customers (name, phone) VALUES ($1, $2) ON CONFLICT DO NOTHING
	`, name, phone)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(resutl.RowsAffected())

	//Update
	id_1 := 1
	newName := "Vasiliy"
	resutl, err = db.ExecContext(ctx, `
	UPDATE customers SET name = $2 WHERE id = $1
	`, id_1, newName)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(resutl.LastInsertId())

	//Select with 1 Row
	// err = db.QueryRowContext(ctx, `
	//  SELECT id, name, phone, active, created FROM customers WHERE id = 1;
	// `).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)

	customer := &types.Customer{}
	id := 1
	newPhone := "+992000000099"
	err = db.QueryRowContext(ctx, `
	  UPDATE customers SET phone = $2 WHERE id = $1 RETURNIGN id, name, phone, active, created
	`, id, newPhone).Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Active, &customer.Created)

	//Обработка пустого результата
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("No Rows")
		return
	}

	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("%v", customer)

	//Select with many Rows
	//создаём слайс для зранения результатов
	items := make([]*types.Customer, 0)
	//делаем запрос
	rows, err := db.QueryContext(ctx, `
	  SELECT id, name, phone, active, created FROM customers WHERE active
	`)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Println(cerr)
		}
	}()
	// rows.Next() - возвращает true до тех пор пока дальше есть строки
	for rows.Next() {
		item := &types.Customer{}
		err = rows.Scan(&item.ID, &item.Name, &item.Phone, &item.Active, &item.Created)
		if err != nil {
			log.Println(err)
			return
		}
		items = append(items, item)
	}
	// и в конце нужно проверять общую сумму
	err = rows.Err()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(items)
*/
