package middleware

import (
	"fmt"

	"draw/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

func AuthMiddleware(authKeys jwk.Set) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := auth.UserFromToken(c.Request, authKeys)
		if err != nil {
			fmt.Println("Auth error:", err)
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Set("userId", userId)
		c.Next()
	}
}
