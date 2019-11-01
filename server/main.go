package main

import (
	"github.com/go-redis/redis/v7"
	"log"
	"net/http"
)

type MyHttpApp struct {
	index       indexHandler
	healthcheck HealthcheckHandler
	jsbundle    indexHandler
	createJob   JobCreateHandler
	listJob     JobListHandler
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

	/*
		read in config and establish connection to persistence layer
	*/
	log.Printf("Reading config from serverconfig.yaml")
	config, configReadErr := ReadConfig("config/serverconfig.yaml")
	log.Print("Done.")

	if configReadErr != nil {
		log.Fatal("No configuration, can't continue")
	}

	redisClient, redisErr := SetupRedis(config)
	if redisErr != nil {
		log.Fatal("Could not connect to redis")
	}

	/*
		configure the elements that each handler requires
	*/
	app.index.filePath = "public/index.html"
	app.index.contentType = "text/html"
	app.healthcheck.redisClient = redisClient
	app.jsbundle.filePath = "public/js/bundle.js"
	app.jsbundle.contentType = "application/javascript"
	app.createJob.redisClient = redisClient
	app.listJob.redisClient = redisClient

	/*
		register each handler to the server
	*/
	http.Handle("/", app.index)
	http.Handle("/healthcheck", app.healthcheck)
	http.Handle("/static/js/bundle.js", app.jsbundle)
	http.Handle("/api/job/start", app.createJob)
	http.Handle("/api/job/list", app.listJob)

	/*
		kick off the server
	*/
	log.Printf("Started HTTP server on port 9000.")
	startServerErr := http.ListenAndServe(":9000", nil)

	if startServerErr != nil {
		log.Fatal(startServerErr)
	}
}
