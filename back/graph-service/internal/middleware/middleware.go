package middleware

import (
	"context"
	"net/http"

	"github.com/0x726f6f6b6965/my-blog/graph-service/internal/client"
	"github.com/0x726f6f6b6965/my-blog/graph-service/internal/utils"
	"github.com/0x726f6f6b6965/my-blog/lib/checker"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// create a context key
type contextKey struct {
	name string
}

// create a context key for user data
var userCtxKey = &contextKey{"user"}

// NewMiddleware returns a middleware for authentication
func NewMiddleware(rds *redis.Client, logger *zap.Logger) func(http.Handler) http.Handler {
	// return handler that acts as a middleware
	return func(next http.Handler) http.Handler {
		// return handler function
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get header data from Authorization header
			var header string = r.Header.Get("Authorization")

			// if header data is empty
			// continue to serve HTTP
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}
			secret, _ := client.GetSecret(context.Background(), rds)

			// get the JWT token from the header
			tokenData, err := utils.CheckToken(r, secret)

			// if the JWT token is invalid, return an error
			// the next request cannot be proceed
			if err != nil {
				logger.Error("check token failed", zap.Error(err))
				http.Error(w, "invalid token", http.StatusForbidden)
				return
			}
			token, _ := client.GetToken(context.Background(), tokenData.UserId, rds)

			// if a user is not found, return an error
			// the next request cannot be proceed
			if checker.IsEmpty(token) {
				http.Error(w, "user not found", http.StatusForbidden)
				return
			}

			// create a context with value
			// the context value is user data
			ctx := context.WithValue(r.Context(), userCtxKey, tokenData.UserId)

			// add context to the request object
			r = r.WithContext(ctx)
			// continue to serve HTTP
			next.ServeHTTP(w, r)
		})
	}
}

func ForContext(ctx context.Context) string {
	raw, _ := ctx.Value(userCtxKey).(string)
	return raw
}
