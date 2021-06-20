package middleware

import (
	"context"
	"net/http"
)

type HasAndRoleFunc func(ctx context.Context, roles ...string) bool

func CheckRole(hasAnyRoleFunc HasAndRoleFunc, roles ...string) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if !hasAnyRoleFunc(request.Context(), roles...) {
				http.Error(writer, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			handler.ServeHTTP(writer, request)
		})
	}
}
