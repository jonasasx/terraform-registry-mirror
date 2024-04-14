package main

import (
	"github.com/eko/gocache/store/go_cache/v4"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"net/http"
	"terraform-registry-mirror/internal/server"
	"time"
)

func main() {
	router := gin.Default()
	cacheClient := cache.New(365*24*time.Hour, 365*24*time.Hour)
	cacheStore := go_cache.NewGoCache(cacheClient)
	serverInstance := server.NewServer(cacheStore)

	router.GET("/:hostname/:namespace/:pkg/index.json", serverInstance.Index)
	router.GET("/:hostname/:namespace/:pkg/:version.json", serverInstance.Version)
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, "OK")
	})
	router.StaticFile("/", "public/index.html")

	err := router.Run("0.0.0.0:8080")
	if err != nil {
		return
	}
}
