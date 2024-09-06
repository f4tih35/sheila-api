package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"strings"
)

type Server struct {
	config      *Config
	redisClient *RedisClient
}

func NewServer(cfg *Config, rdb *RedisClient) *Server {
	return &Server{
		config:      cfg,
		redisClient: rdb,
	}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", s.config.TCPAddress)
	if err != nil {
		log.Fatalf("Failed to start TCP server: %v", err)
	}
	defer listener.Close()

	log.Printf("TCP server listening on %s...", s.config.TCPAddress)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().String()
	ip := extractIP(remoteAddr)
	log.Printf("New connection from: %s", ip)

	if err := s.redisClient.AddIP(ip); err != nil {
		log.Printf("Error adding IP: %v", err)
		return
	}

	defer func() {
		if err := s.redisClient.RemoveIP(ip); err != nil {
			log.Printf("Error removing IP: %v", err)
		}
		log.Printf("Connection closed: %s", ip)
	}()

	otherIPs, err := s.redisClient.GetAllIPs()
	if err != nil {
		log.Printf("Error retrieving IPs: %v", err)
		return
	}

	filteredIPs := filterIPs(otherIPs, ip)

	jsonData, err := json.Marshal(filteredIPs)
	if err != nil {
		log.Printf("Error creating JSON: %v", err)
		return
	}

	_, err = conn.Write(jsonData)
	if err != nil {
		log.Printf("Error sending data: %v", err)
		return
	}

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.ToLower(text) == "exit" {
			break
		}
	}
}

func extractIP(remoteAddr string) string {
	parts := strings.Split(remoteAddr, ":")
	if len(parts) < 1 {
		return remoteAddr
	}
	return parts[0]
}

func filterIPs(ips []string, excludeIP string) []string {
	var filtered []string
	for _, ip := range ips {
		if ip != excludeIP {
			filtered = append(filtered, ip)
		}
	}
	return filtered
}
