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
	"github.com/SardorMS/CRUD/pkg/managers"
)

func main() {

	host := "127.0.0.1"
	port := "9999"
	//user:login@host:port/db
	dsn := "postgres://app:123@192.168.99.100:5432/db"

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
		managers.NewService,
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
