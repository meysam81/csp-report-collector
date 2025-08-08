package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(ctx context.Context, c *Config) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.RedisHost, c.RedisPort),
		Password: c.RedisPassword,
		DB:       c.RedisDB,
	}

	if c.RedisSSLRequired {
		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
		}
	}

	redisClient := redis.NewClient(opts)
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		ctxS, cancelS := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelS()
		_, err := redisClient.Shutdown(ctxS).Result()
		if err != nil {
			fmt.Printf("failed shutting down the redis: %s", err)
		}
	}()

	return redisClient, nil
}
