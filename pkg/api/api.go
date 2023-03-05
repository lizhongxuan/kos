package api

import (
	"github.com/gin-gonic/gin"
	"kos/pkg/service"
)

func RegisterApi(addr string) {
	router := gin.Default()
	router.GET("/objects/get", service.GetObject)
	router.PUT("/objects/put", service.PutObject)
	router.Run(addr)
}
