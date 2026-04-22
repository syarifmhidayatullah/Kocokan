package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/project/kocokan/internal/service"
	"github.com/project/kocokan/pkg/response"
)

const UserKey = "user_id"

func Auth(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			// try cookie
			cookie, err := c.Cookie("koco_token")
			if err != nil || cookie == "" {
				response.Error(c, 401, "Unauthorized")
				c.Abort()
				return
			}
			header = "Bearer " + cookie
		}
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, 401, "Token tidak valid")
			c.Abort()
			return
		}
		user, err := authSvc.ValidateToken(parts[1])
		if err != nil {
			response.Error(c, 401, "Token tidak valid atau kadaluwarsa")
			c.Abort()
			return
		}
		c.Set(UserKey, user.ID)
		c.Next()
	}
}

func UserID(c *gin.Context) uint {
	v, _ := c.Get(UserKey)
	id, _ := v.(uint)
	return id
}
