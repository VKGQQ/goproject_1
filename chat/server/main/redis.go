package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func initRedis(address string, password string, db int, poolSize int) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:            address,
		Password:        password,
		DB:              db,
		PoolSize:        poolSize, // 连接池最大连接数
		MinIdleConns:    3,        // 最小空闲连接数
		DialTimeout:     5 * time.Second,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
		ConnMaxIdleTime: 5 * time.Minute,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		fmt.Println("redis 连接失败：", err)
		return
	}
	fmt.Println("redis 连接成功")
}
