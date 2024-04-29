package router

import (
	"github.com/0x726f6f6b6965/tiny-url-go/cmd/api"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine, short *api.ShortenAPI) {
	server.POST("/shorten", short.Shorten)
	server.GET("/:shorten", short.RedirectURL)
	server.DELETE("/shorten", short.DeleteURL)
	server.PATCH("/shorten", short.UpdateURL)

}
