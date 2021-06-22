package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
)

const (
	MANAGER = "MANAGER"
	ADMIN   = "ADMIN"
)

var ErrNoAuthentication = errors.New("no authentication")

// A variable that will be the key by which the value will be added.
var authenticationContextKey = &contextKey{"authentication context"}

// Non-exportable type.
type contextKey struct {
	name string
}

func (c *contextKey) String() string {
	return c.name
}

type IDFunc func(ctx context.Context, token string) (int64, error)

// Authenticate - authentication procedure
func Authenticate(idFunc IDFunc) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			token := request.Header.Get("Authorization")

			id, err := idFunc(request.Context(), token)
			if err != nil {
				log.Println(err, "Not Authorized")
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			// give a value by a key.
			ctx := context.WithValue(request.Context(), authenticationContextKey, id)
			request = request.WithContext(ctx)

			handler.ServeHTTP(writer, request)
		})
	}
}

// Athuntecation - helper function, to extract value from context.
func Authentication(ctx context.Context) (int64, error) {
	if value, ok := ctx.Value(authenticationContextKey).(int64); ok {
		return value, nil
	}
	return 0, ErrNoAuthentication
}
