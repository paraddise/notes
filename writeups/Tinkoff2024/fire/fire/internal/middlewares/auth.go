package middlewares

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user_id")
	if user == nil {
		c.Redirect(http.StatusSeeOther, "/signin")
		c.Abort()
		return
	}
	c.Next()
}
