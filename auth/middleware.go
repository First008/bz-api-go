package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO bunu project globali olarak kullanabilirmiyiz ?

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := isTokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, ` {"message": "unauthorised"}`)
			c.Abort()
			return
		}
		c.Next()
	}
}
