package isulog

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pirosiki197/isulog/internal"
)

func Gin() gin.HandlerFunc {
	return GinWithConfig(DefaultConfig)
}

func GinWithConfig(config Config) gin.HandlerFunc {
	recorder := internal.NewRecorder(config.Filename)

	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		r := internal.Record{
			Path:         ctx.FullPath(),
			Method:       ctx.Request.Method,
			StatusCode:   ctx.Writer.Status(),
			ResponseTime: time.Since(start),
		}
		if err := recorder.Save(r); err != nil {
			log.Printf("[ERROR] Failed to record. isulog: %s", err.Error())
		}
	}
}
