package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	g := gin.Default()

	admin := g.Group("/api/admin")
	{
		admin.GET("/", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "current time: %v", time.Now().Format(time.RFC3339))
		})
	}

	if err := g.Run(":3000"); err != nil {
		panic(err)
	}
}
