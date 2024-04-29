package main

import (
	"net/http"

	"github.com/0x726f6f6b6965/tiny-url-go/cmd/api"
	"github.com/0x726f6f6b6965/tiny-url-go/cmd/api/router"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func initServer(env string, ser *api.ShortenAPI) *gin.Engine {
	gin.SetMode(func() string {
		if env == "dev" {
			return gin.DebugMode
		}
		return gin.ReleaseMode
	}())
	engine := gin.New()
	engine.Use(cors.Default())
	engine.Use(gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "Service internal exception!",
		})
	}))
	router.RegisterRoutes(engine, ser)
	return engine
}
