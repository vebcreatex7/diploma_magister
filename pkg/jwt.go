package pkg

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/vebcreatex7/diploma_magister/internal/api/response"
	"net/http"
	"os"
	"strconv"
	"time"
)

var privateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))

// generate JWT token
func GenerateJWT(user response.User) (string, error) {
	tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":  user.UID,
		"role": user.Role,
		"iat":  time.Now().Unix(),
		"eat":  time.Now().Add(time.Second * time.Duration(tokenTTL)).Unix(),
	})
	return token.SignedString(privateKey)
}

// validate JWT token
func ValidateJWT(r *http.Request) error {
	token, err := getToken(r)
	if err != nil {
		return err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}
	return errors.New("invalid token provided")
}

// validate Admin role
func ValidateAdminRoleJWT(r *http.Request) error {
	token, err := getToken(r)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	userRole := claims["role"].(string)
	if ok && token.Valid && userRole == "admin" {
		return nil
	}
	return errors.New("invalid admin token provided")
}

// validate Customer role
func ValidateScientistRoleJWT(r *http.Request) error {
	token, err := getToken(r)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	userRole := claims["role"].(string)
	if ok && token.Valid && userRole == "scientist" {
		return nil
	}
	return errors.New("invalid scientist token provided")
}

func getToken(r *http.Request) (*jwt.Token, error) {
	for _, c := range r.Cookies() {
		if c.Name == "jwt" {
			token, err := jwt.Parse(c.Value, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				return privateKey, nil
			})
			if err != nil {
				return nil, err
			}

			return token, nil
		}
	}

	return nil, fmt.Errorf("no token")
}

func ValidateAdminJWTCookies(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		userRole := claims["role"].(string)
		if ok && token.Valid && userRole == "admin" {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	})
}

func ValidateScientistJWTCookies(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		userRole := claims["role"].(string)
		if ok && token.Valid && userRole == "scientist" {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	})
}

func ValidateEngineerJWTCookies(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		userRole := claims["role"].(string)
		if ok && token.Valid && userRole == "engineer" {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	})
}

func ValidateLaboratorianJWTCookies(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		userRole := claims["role"].(string)
		if ok && token.Valid && userRole == "laboratorian" {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	})
}

func GetUIDFromJWT(r *http.Request) (string, error) {
	token, err := getToken(r)
	if err != nil {
		return "", fmt.Errorf("getting token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	uid := claims["uid"].(string)

	if !ok || uid == "" {
		return "", fmt.Errorf("uid not found")
	}

	return uid, nil
}
