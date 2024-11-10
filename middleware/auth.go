package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5")

// AuthMiddleware untuk memvalidasi token JWT
func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mengambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Memeriksa format header (harus 'Bearer <token>')
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Mem-parse token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Memeriksa metode signing token
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secretKey), nil
		})

		// Jika token tidak valid
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: invalid token"})
			c.Abort()
			return
		}

		// Mengambil klaim dari token dan menyimpannya di context
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if authID, ok := claims["auth_id"].(float64); ok {
				c.Set("auth_id", int64(authID))
			}
			if accountID, ok := claims["account_id"].(float64); ok {
				c.Set("account_id", int64(accountID))
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: invalid claims"})
			c.Abort()
			return
		}

		// Jika token valid, lanjutkan ke handler berikutnya
		c.Next()
	}
}
