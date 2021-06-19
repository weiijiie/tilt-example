package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/kelseyhightower/envconfig"
)

type (
	SetRequest struct {
		Value string `json:"value" binding:"required"`
	}

	Config struct {
		Port          int    `default:"3000"`
		RedisAddr     string `required:"true" split_words:"true"`
		RedisPassword string `required:"true" split_words:"true"`
	}
)

func ping(store Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := store.Ping(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"msg": result})
		return
	}
}

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
	var conf Config
	err := envconfig.Process("", &conf)
	if err != nil {
		panic(err.Error())
	}

	store := NewStore(conf.RedisAddr, conf.RedisPassword)
	go func() {
		for range time.Tick(5 * time.Second) {
			ctx := context.Background()
			res, err := store.Ping(ctx)
			if err != nil {
				log.Printf("heartbeat failed, err: %s", err.Error())
			}
			log.Printf("heartbeat: %s...", res)
		}
	}()

	router := gin.Default()

	router.GET("/", ping(store))
	router.GET("/:key", get(store))
	router.POST("/:key", set(store))

	if err := router.Run(":" + strconv.Itoa(conf.Port)); err != nil {
		log.Fatal(err)
	}
}
