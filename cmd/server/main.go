package main

import (
	"github.com/eko/gocache/store/go_cache/v4"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"net/http"
	"strings"
	"terraform-registry-mirror/internal/server"
	"time"
)

func main() {
	router := gin.New()

	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/healthz"},
	}))
	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		c.Next() // Выполняем обработку запроса

		// Получаем текущий Content-Type
		contentType := c.Writer.Header().Get("Content-Type")

		// Проверяем, содержит ли Content-Type "application/json"
		if strings.HasPrefix(contentType, "application/json") {
			// Устанавливаем Content-Type только как "application/json"
			c.Writer.Header().Set("Content-Type", "application/json")
		}
	})

	cacheClient := cache.New(365*24*time.Hour, 365*24*time.Hour)
	cacheStore := go_cache.NewGoCache(cacheClient)
	serverInstance := server.NewServer(cacheStore)

	router.GET("/:hostname/:namespace/:pkg/index.json", serverInstance.Index)
	router.GET("/:hostname/:namespace/:pkg/:version.json", serverInstance.Version)
	router.GET("/healthz", func(c *gin.Context) {
		c.AsciiJSON(http.StatusOK, "OK")
	})
	router.StaticFile("/", "public/index.html")

	err := router.Run("0.0.0.0:8080")
	if err != nil {
		return
	}
}
