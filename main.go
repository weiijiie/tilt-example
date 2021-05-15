package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"os"
)

type (
	SetRequest struct {
		Value string `json:"value" binding:"required"`
	}
)

func get(store Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Param("key")
		result, err := store.Get(c, key)
		if err == redis.Nil {
			c.JSON(http.StatusNotFound, nil)
			return
		} else if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"value": result})
	}
}

func set(store Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SetRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		key := c.Param("key")
		result, err := store.Set(c, key, req.Value)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		}

		log.Printf("set successful, key: %s, val: %s", key, result)
		c.JSON(http.StatusOK, gin.H{"msg": fmt.Sprintf("%s set for key %s", result, key)})
	}
}

func main() {
	var port string
	if isDev() {
		port = os.Getenv("PORT")
	} else {
		port = "80"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if port == "" || redisAddr == "" {
		panic("config not set properly")
	}

	store := NewStore(redisAddr)

	router := gin.Default()

	router.GET("/:key", get(store))
	router.POST("/:key", set(store))

	if port == "" {
		panic("`PORT` env var not set")
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func isDev() bool {
	return os.Getenv("IS_DEV") == "true"
}
