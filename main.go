package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var cache = make(map[string]string)
var m sync.Mutex

func getHandler(c *gin.Context) {
	m.Lock()
	defer m.Unlock()

	key := c.Param("key")

	value, exists := cache[key]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
}

func upsertHandler(c *gin.Context) {
	m.Lock()
	defer m.Unlock()

	key := c.Query("key")

	value, exists := cache[key]
	if exists {
		c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
		return
	}

	cache[key] = c.Query("value")
	c.JSON(http.StatusOK, gin.H{"status": "ok", "key": key})
}

func main() {
	r := gin.Default()

	r.POST("/set", upsertHandler)
	r.GET("/get/:key", getHandler)

	r.Run(":8080")
}
