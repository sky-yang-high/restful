package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		log.Println("Request", ctx.Request.Method, ctx.Request.URL.Path, "took", time.Since(start))
	}
}

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Recovery from panic:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"err": err})
			}
		}()
		ctx.Next()
	}
}

func Authrization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, psw, ok := ctx.Request.BasicAuth()
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}
		if !CheckUser(user, psw) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func CheckUser(user, psw string) bool {
	// TODO: check user and password
	return true
}
