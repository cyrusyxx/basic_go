package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
}

func (m LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		gob.Register(time.Now())

		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			return
		}
		sess := sessions.Default(ctx)
		uid := sess.Get("userId")
		if uid == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		const updateTimeKey = "update_time"
		utime, ok := sess.Get(updateTimeKey).(time.Time)
		if !ok || time.Now().Sub(utime) > time.Minute {
			sess.Set(updateTimeKey, time.Now())
			sess.Set("userId", uid)
			err := sess.Save()
			if err != nil {

			}
		}
	}
}
