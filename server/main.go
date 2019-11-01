package main

import (
	"github.com/go-redis/redis/v7"
	"log"
	"net/http"
)

type MyHttpApp struct {
	index    indexHandler
	jsbundle indexHandler
}

func SetupRedis(config *Config) (*redis.Client, error) {
	log.Printf("Connecting to Redis on %s", config.Redis.Address)
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Address,
		Password: config.Redis.Password,
		DB:       config.Redis.DBNum,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Printf("Could not contact Redis: %s", err)
		return nil, err
	}
	log.Printf("Done.")
	return client, nil
}

func main() {
	var app MyHttpApp

	log.Printf("Reading config from serverconfig.yaml")
	config, configReadErr := ReadConfig("serverconfig.yaml")
	log.Print("Done.")

	if configReadErr != nil {
		log.Fatal("No configuration, can't continue")
	}

	_, redisErr := SetupRedis(config)
	if redisErr != nil {
		log.Fatal("Could not connect to redis")
	}

	app.index.filePath = "public/index.html"
	app.index.contentType = "text/html"
	app.jsbundle.filePath = "public/js/bundle.js"
	app.jsbundle.contentType = "application/javascript"

	http.Handle("/", app.index)
	http.Handle("/static/js/bundle.js", app.jsbundle)

	log.Printf("Started HTTP server on port 9000.")
	startServerErr := http.ListenAndServe(":9000", nil)

	if startServerErr != nil {
		log.Fatal(startServerErr)
	}
}
