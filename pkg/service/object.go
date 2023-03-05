package service

import (
	"io"
	"net/http"
	"os"
	"strings"
)
import "github.com/gin-gonic/gin"

func GetObject(c *gin.Context) {
	f, err := os.Create(os.Getenv("STORAGE_ROOT") + "/objects/" + strings.Split(c.Request.URL.EscapedPath(), "/")[2])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "not find",
		})
		return
	}
	defer f.Close()
	_, err = io.Copy(c.Writer, f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad copy",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
	return
}
func PutObject(c *gin.Context) {
	f, err := os.Create(os.Getenv("STORAGE_ROOT") + "/objects/" + strings.Split(c.Request.URL.EscapedPath(), "/")[2])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "not find",
		})
		return
	}
	defer f.Close()
	_, err = io.Copy(f, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad copy",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
	return
}
