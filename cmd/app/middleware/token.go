package middleware

import (
	"context"
	"errors"
	"net/http"
)

var ErrNoAuthentication = errors.New("no authentication")

//Переменную этого типа, которая и будет ключом по которому будет
//класться значение
var authenticationContextKey = &contextKey{"authentication context"}

// Неэкспортируемый тип
type contextKey struct {
	name string
}

func (c *contextKey) String() string {
	return c.name
}

type IDFunc func(ctx context.Context, token string) (int64, error)

// Authenticate - ...
func Authenticate(idFunc IDFunc) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			token := request.Header.Get("Authorization")

			id, err := idFunc(request.Context(), token)
			if err != nil {
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			// Кладём по этому ключу значение
			ctx := context.WithValue(request.Context(), authenticationContextKey, id)
			request = request.WithContext(ctx)

			handler.ServeHTTP(writer, request)
		})
	}
}

// Athuntecation - функцию helper*, чтобы доставать значение из контекста
func Authentication(ctx context.Context) (int64, error) {
	if value, ok := ctx.Value(authenticationContextKey).(int64); ok {
		return value, nil
	}
	return 0, ErrNoAuthentication
}
