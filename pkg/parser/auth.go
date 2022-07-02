package parser

import (
	"context"
	"errors"
	"time"

	"github.com/domino14/tshparser/rpc/proto"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const TokenExpire = 60 * 24 * time.Hour

type AuthService struct {
	store     *SqliteStore
	secretKey string
}

func NewAuthService(store *SqliteStore, secretKey string) *AuthService {
	return &AuthService{store: store, secretKey: secretKey}
}

func (a *AuthService) RegisterUser(ctx context.Context, req *proto.NewUserRequest) (*proto.NewUserResponse, error) {
	if len(req.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}
	if len(req.Email) < 3 {
		return nil, errors.New("email not provided")
	}
	hash, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	err = a.store.NewUser(ctx, req.Email, hash)
	if err != nil {
		return nil, err
	}
	return &proto.NewUserResponse{}, nil
}

func hashPassword(password string) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Hash password with Bcrypt's min cost
	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

func doPasswordsMatch(hashedPassword, currPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(currPassword))
	return err == nil
}

func (a *AuthService) GetJWT(ctx context.Context, req *proto.JWTRequest) (*proto.JWTResponse, error) {
	user, err := a.store.GetUser(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if !doPasswordsMatch(user.PasswordHash, req.Password) {
		return nil, errors.New("passwords do not match")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": req.Email,
		"adm": user.IsAdmin,
		"exp": time.Now().Add(TokenExpire).Unix(),
	})
	tokenString, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return nil, err
	}
	SetDefaultCookie(ctx, tokenString)
	return &proto.JWTResponse{Token: tokenString}, nil
}

// A test endpoint just to verify a token.
func (a *AuthService) WhoAmI(ctx context.Context, req *proto.WhoAmIRequest) (*proto.User, error) {
	return userFromHTTPHeader(ctx, a.secretKey)
}

func userFromHTTPHeader(ctx context.Context, secretKey string) (*proto.User, error) {
	tokenString, err := GetPassedInJWT(ctx)
	if err != nil {
		return nil, err
	}
	return userFromToken(tokenString, secretKey)
}

func userFromToken(tokenString, secretKey string) (*proto.User, error) {

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	user := &proto.User{}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if email, ok := claims["sub"].(string); !ok {
			return nil, errors.New("email malformed")
		} else {
			user.Email = email
		}
		if adm, ok := claims["adm"].(bool); !ok {
			return nil, errors.New("adm malformed")
		} else {
			user.IsAdmin = adm
		}
	} else {
		return nil, errors.New("invalid token")
	}
	return user, nil
}
