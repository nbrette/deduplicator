package main

import (
	"deduplicator/internal/controllers"
	"deduplicator/internal/services"
	"deduplicator/internal/services/router"
	"log"
	"net/http"
	"os"
)

func main() {
	redisHost := os.Getenv("REDIS_URL")
	log.Printf("HOST: %s", redisHost)
	redisService, err := services.NewRedisService(redisHost, "", 0)
	if err != nil {
		log.Panicf("Can't connect to Redis DB : %s", err)
	}
	deduplicatorController := controllers.NewDeduplicatorController(services.NewDeduplicatorService(redisService), services.NewPubSubService(), services.NewHTTPService())
	log.Printf("Running")
	log.Fatal(http.ListenAndServe(":8080", router.InitializeRouter(deduplicatorController)))
}
