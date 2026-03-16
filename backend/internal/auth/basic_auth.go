package auth

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	DB *sql.DB
}

func parseBasicAuth(header string) (string, string, error) {
	if header == "" {
		return "", "", errors.New("missing authorization header")
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Basic" {
		return "", "", errors.New("invalid authorization scheme")
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", err
	}

	creds := strings.SplitN(string(decoded), ":", 2)
	if len(creds) != 2 {
		return "", "", errors.New("invalid basic auth format")
	}

	return creds[0], creds[1], nil
}

func (a *AuthService) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, err := parseBasicAuth(c.GetHeader("Authorization"))
		if err != nil {
			c.Header("WWW-Authenticate", `Basic realm="pet-store"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}

		var userID, dbUsername, passwordHash, role string
		var storeID sql.NullString

		err = a.DB.QueryRow(`
			SELECT id, username, password_hash, role, store_id
			FROM users
			WHERE username = $1
		`, username).Scan(&userID, &dbUsername, &passwordHash, &role, &storeID)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		var storeIDPtr *string
		if storeID.Valid {
			storeIDPtr = &storeID.String
		}

		ctxUser := &UserContext{
			UserID:   userID,
			Username: dbUsername,
			Role:     role,
			StoreID:  storeIDPtr,
		}

		ctx := WithUser(c.Request.Context(), ctxUser)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
