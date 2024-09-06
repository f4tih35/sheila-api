package main

func main() {
	cfg := LoadConfig()
	redisClient := NewRedisClient(cfg)
	server := NewServer(cfg, redisClient)
	server.Start()
}
