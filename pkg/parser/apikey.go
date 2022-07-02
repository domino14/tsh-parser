package parser

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/twitchtv/twirp"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ctxkey string

const ApiKeyHeader = "Authorization"
const apikeykey ctxkey = "apikey"
const rwkey ctxkey = "responsewriter"

// ExposeResponseWriterMiddleware configures an http.Handler (like any Twirp server)
// to place the responseWriter in its context. This should enable
// setting cookies with the setCookie function.
func ExposeResponseWriterMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, rwkey, w)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

func SetDefaultCookie(ctx context.Context, value string) error {
	cookie := &http.Cookie{
		Name:     "session",
		Value:    value,
		Expires:  time.Now().Add(60 * 24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	}
	log.Debug().Msgf("setting cookie %v", cookie)
	return setCookie(ctx, cookie)
}

func setCookie(ctx context.Context, cookie *http.Cookie) error {
	w, ok := ctx.Value(rwkey).(http.ResponseWriter)
	if !ok {
		return errors.New("unable to get ResponseWriter from context, middleware might not be set up correctly")
	}
	http.SetCookie(w, cookie)
	return nil
}

// JWTMiddleware creates a middleware to fetch an API key from
// a header and store it in a context key.
func JWTMiddleware(h http.Handler) http.Handler {
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

// AuthenticationMiddleware is an auth middleware to fetch an API
// key from the Cookie header and store it in a context key.
func AuthenticationMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := zerolog.Ctx(ctx)
		var err error
		sessionCookie, err := r.Cookie("session")
		if err != nil {
			if err != http.ErrNoCookie {
				log.Err(err).Msg("error-getting-new-cookie")
			}
			log.Debug().Msg("no cookie")
			h.ServeHTTP(w, r)
			return
		}
		ctx = context.WithValue(ctx, apikeykey, sessionCookie.Value)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
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
