package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisClient(cfg *Config) *RedisClient {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return &RedisClient{
		client: rdb,
		ctx:    ctx,
	}
}

func (r *RedisClient) AddIP(ip string) error {
	return r.client.SAdd(r.ctx, "connected_ips", ip).Err()
}

func (r *RedisClient) RemoveIP(ip string) error {
	return r.client.SRem(r.ctx, "connected_ips", ip).Err()
}

func (r *RedisClient) GetAllIPs() ([]string, error) {
	return r.client.SMembers(r.ctx, "connected_ips").Result()
}
