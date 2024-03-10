package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// SetupRoutes connects the HTTP API endpoints to the handlers
func SetupRoutes() *gin.Engine {
	r := gin.Default()

	r.Use(cors.Default())

	r.Use(static.Serve("/orders", static.LocalFile("./static", true)))
	r.Static("/css", "public/css")
	r.Static("/js", "public/js")

	r.GET("/api/order/list", List)
	r.GET("/api/order/import", Import)

	return r
}
