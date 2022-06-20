package parser

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/twitchtv/twirp"

	"github.com/rs/zerolog"
)

type ctxkey string

const ApiKeyHeader = "Authorization"
const apikeykey ctxkey = "apikey"

// JWTMiddlewareGenerator creates a middleware to fetch an API key from
// a header and store it in a context key.
func JWTMiddlewareGenerator() (mw func(http.Handler) http.Handler) {
	mw = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := zerolog.Ctx(ctx)

			apikey := r.Header[ApiKeyHeader]
			if len(apikey) > 1 {
				log.Error().Msg("apikey formatted incorrectly")
				// Serve, unauthenticated
				h.ServeHTTP(w, r)
				return
			} else if len(apikey) == 0 {
				h.ServeHTTP(w, r)
				return
			}
			_, jwt, found := strings.Cut(apikey[0], "Bearer ")
			if !found {
				h.ServeHTTP(w, r)
			}
			// Otherwise, a JWT was provided. Store it in the context.
			ctx = context.WithValue(ctx, apikeykey, jwt)
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
	return
}

// GetPassedInJWT works with JWTMiddlewareGenerator to return an API key in the
// passed-in context.
func GetPassedInJWT(ctx context.Context) (string, error) {
	apikey := ctx.Value(apikeykey)
	if apikey == nil {
		return "", twirp.NewError(twirp.Unauthenticated, "authentication required")
	}
	a, ok := apikey.(string)
	if !ok {
		return "", twirp.InternalErrorWith(errors.New("unexpected error with apikey type inference"))
	}
	return a, nil
}
