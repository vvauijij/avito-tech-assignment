package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/vvauijij/avito-tech-assignment/server/internal/env"
	"github.com/vvauijij/avito-tech-assignment/server/internal/handler"
	"github.com/vvauijij/avito-tech-assignment/server/internal/mongo"
	"github.com/vvauijij/avito-tech-assignment/server/internal/redis"
	"github.com/vvauijij/avito-tech-assignment/server/internal/token"
)

var (
	ENV = os.Getenv("ENV")

	serverPort = flag.Int("server", 8080, "HTTP server port")
	redisURI   = flag.String("redis", "", "Redis URI")
	mongoURI   = flag.String("mongo", "", "Mongo URI")
	publicFile = flag.String("public", "", "path to JWT public key file")
)

func main() {
	flag.Parse()
	if serverPort == nil {
		fmt.Fprintln(os.Stderr, "HTTP server port is required")
		os.Exit(1)
	}
	if redisURI == nil || *redisURI == "" {
		fmt.Fprintln(os.Stderr, "Redis URI is required")
		os.Exit(1)
	}
	if mongoURI == nil || *mongoURI == "" {
		fmt.Fprintln(os.Stderr, "Mongo URI is required")
		os.Exit(1)
	}
	if publicFile == nil || *publicFile == "" {
		fmt.Fprintln(os.Stderr, "path to JWT public key file is required")
		os.Exit(1)
	}

	redisClient := redis.NewClient(*redisURI, 5*time.Minute)
	mongoClient := mongo.NewClient(*mongoURI)
	tokenClient := token.NewClient(*publicFile)

	httpHandler := handler.NewHTTPHandler(tokenClient, redisClient, mongoClient)
	http.HandleFunc("/user_banner", httpHandler.UserBanner)
	http.HandleFunc("/banner", httpHandler.AdminBanner)
	http.HandleFunc("/banner/{id}", httpHandler.AdminBannerWithID)
	if env.IsTestENV() {
		http.HandleFunc("/test_clean_up", httpHandler.TestCleanUp)
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", *serverPort), nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
