package world_cup_api

import (
	"fmt"
	"github.com/amitizle/telegram-world-cup-bot/internal/http_client"
	"github.com/go-redis/redis"
	"log"
	"time"
)

var (
	healthcheckInterval     = 2 * time.Second
	todayMatchesInterval    = 1 * time.Minute
	currentMatchesInterval  = 1 * time.Minute
	tomorrowMatchesInterval = 1 * time.Hour
)

func StartPolling(redisHost string, redisPort int) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", redisHost, redisPort),
		DB:   0, // default db
	})
	go heartbeat(redisClient)
	go tomorrowMatches(redisClient)
	go todayMatches(redisClient)
	go currentMatches(redisClient)
}

// TODO configurable hard coded "today_matches"
func tomorrowMatches(redisClient *redis.Client) {
	for {
		httpClient, err := world_cup_http_client.New("")
		if err != nil {
			log.Printf("error while creating new HTTP client: %v", err)
		}
		response, err := httpClient.Get("/matches/tomorrow", map[string]string{})
		if err != nil {
			log.Printf("error while HTTP getting /matches/tomorrow: %v", err)
		}
		err = redisClient.Set("tomorrow_matches", response.Body, 0).Err()
		if err != nil {
			log.Printf("error while setting key in Redis: %v", err)
		}
		time.Sleep(tomorrowMatchesInterval)
	}
}

// TODO configurable hard coded "today_matches"
func todayMatches(redisClient *redis.Client) {
	for {
		httpClient, err := world_cup_http_client.New("")
		if err != nil {
			log.Printf("error while creating new HTTP client: %v", err)
		}
		response, err := httpClient.Get("/matches/today", map[string]string{})
		if err != nil {
			log.Printf("error while HTTP getting /matches/today: %v", err)
		}
		err = redisClient.Set("today_matches", response.Body, 0).Err()
		if err != nil {
			log.Printf("error while setting key in Redis: %v", err)
		}
		time.Sleep(todayMatchesInterval)
	}
}

// TODO configurable hard coded "current_matches"
func currentMatches(redisClient *redis.Client) {
	for {
		httpClient, err := world_cup_http_client.New("")
		if err != nil {
			log.Printf("error while creating new HTTP client: %v", err)
		}
		response, err := httpClient.Get("/matches/current", map[string]string{})
		if err != nil {
			log.Printf("error while HTTP getting /matches/current: %v", err)
		}
		err = redisClient.Set("current_matches", response.Body, 0).Err()
		if err != nil {
			log.Printf("error while setting key in Redis: %v", err)
		}
		time.Sleep(currentMatchesInterval)
	}
}

func heartbeat(redisClient *redis.Client) {
	for {
		_, err := redisClient.Ping().Result()
		if err != nil {
			log.Printf("redis pong error: %v", err)
		}
		time.Sleep(healthcheckInterval)
	}
}
