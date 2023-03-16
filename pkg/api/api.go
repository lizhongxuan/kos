package api

import (
	"github.com/gin-gonic/gin"
	"kos/pkg/service"
)

func RegisterApi(addr string) {
	router := gin.Default()
	bucketRouter := router.Group("/:bucket")
	bucketRouter.GET("/:objects/get", service.GetObject)
	bucketRouter.PUT("/:objects/put", service.PutObject)

	router.Run(addr)
}
