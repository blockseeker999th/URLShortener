package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/blockseeker999th/URLShortener/internal/config"
	logUtils "github.com/blockseeker999th/URLShortener/internal/utils"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func WithAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := GetTokenFromRequest(r)
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		token, err := validateJWT(tokenString)

		if err != nil {
			slog.Error("Failed to authenticate")
			permissionDenied(w, r)
			return
		}

		if !token.Valid {
			permissionDenied(w, r)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userId := claims["userID"].(string)

		ctx := context.WithValue(r.Context(), "userId", userId)

		handlerFunc(w, r.WithContext(ctx))
	}
}

func permissionDenied(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, logUtils.ResponseWithoutPayload(http.StatusUnauthorized, "permission denied"))
}

func CreateJWT(secret []byte, userId int64) (string, error) {
	const op = "auth.CreateJWT"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(int(userId)),
		"expiresAt": time.Now().Add(time.Hour * 24 * 120).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return tokenString, nil
}

func validateJWT(token string) (*jwt.Token, error) {
	secret := config.MustLoad().JwtSecret
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func HashPassword(password string) (string, error) {
	const op = "auth.HashPassword"

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return string(hash), nil
}

func GetTokenFromRequest(r *http.Request) string {
	tokenHeader := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenHeader != "" {
		return tokenHeader
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}

func CreateAndSetAuthCookie(id int64, w http.ResponseWriter) (string, error) {
	const op = "auth.CreateAndSetAuthCookie"

	secret := []byte(config.MustLoad().JwtSecret)
	token, err := CreateJWT(secret, id)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "Authorization",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 24),
	})

	return token, nil
}
