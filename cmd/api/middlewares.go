package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func setJWTToCookie(c echo.Context, tokenString string) error {
	// Create the cookie with token
	cookie := http.Cookie{
		Name:     "session_token",
		Value:    tokenString,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24),
	}

	http.SetCookie(c.Response().Writer, &cookie)
	return nil
}

func AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_token")
		if err != nil {
			return errorJSON(c.Response().Writer, err, http.StatusUnauthorized)
		}

		tokenString := cookie.Value
		id, err := verifyToken(tokenString)
		if err != nil {
			return errorJSON(c.Response().Writer, err, http.StatusUnauthorized)
		}

		c.Set("id", id)

		return next(c)
	}
}

func verifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if ok != true {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", fmt.Errorf("Invalid expiration time")
	}

	if time.Now().Unix() > int64(exp) {
		return "", fmt.Errorf("Token expired")
	}

	idInt := int(claims["id"].(float64))
	id := strconv.Itoa(idInt)

	return id, nil
}

func RoleRequiredMiddleware(next echo.HandlerFunc, role string) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session_token")
		if err != nil {
			return errorJSON(c.Response().Writer, err, http.StatusUnauthorized)
		}

		cookieValue := cookie.Value

		token, err := jwt.Parse(cookieValue, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}

			return []byte("notsecret"), nil
		})

		claims, _ := token.Claims.(jwt.MapClaims)
		customErr := fmt.Errorf("%s", "Unauthorized!!")
		if claims["role"] != role {
			return errorJSON(c.Response().Writer, customErr, http.StatusUnauthorized)

		}

		return next(c)
	}
}
