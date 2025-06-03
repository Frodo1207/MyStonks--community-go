package middleware

import (
	"MyStonks-go/internal/common/response"
	"MyStonks-go/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if solAddress, ok := validateToken(token); !ok {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse(response.ErrorCodeInvalidToken, []string{}))
			c.Abort()
			return
		} else {
			c.Set("sol_address", solAddress)
			c.Next()
		}
	}
}

func validateToken(token string) (string, bool) {
	const prefix = "Bearer "
	if !strings.HasPrefix(token, prefix) {
		return "", false
	}

	realToken := strings.TrimPrefix(token, prefix)
	claims, err := service.ValidateToken(realToken, service.TokenTypeAccess)
	if err != nil {
		return "", false
	}

	solAddress := claims.WalletAddress

	return solAddress, true
}
