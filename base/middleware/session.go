package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func (m *Middleware) Session() gin.HandlerFunc {
	store := cookie.NewStore([]byte("lkfsxaqws"))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   600,
		HttpOnly: true,
		// 如果是跨域，可能需要加上：
		// SameSite: http.SameSiteNoneMode,
		// Secure:   true,
	})
	return sessions.Sessions("qqlx", store)
}
