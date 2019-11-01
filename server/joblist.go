package main

import (
	"github.com/go-redis/redis/v7"
	"log"
	"net/http"
)

type JobListResponse struct {
	Status string      `json:"status"`
	Jobs   []*JobEntry `json:"jobs"`
	Count  int         `json:"count"`
}

type JobListHandler struct {
	redisClient *redis.Client
}

func (h JobListHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	if AssertHttpMethod(request, w, "GET") == false {
		return
	}

	responses, err := ListJobs(h.redisClient, 100)

	if err != nil {
		log.Printf("Could not list jobs from database: %s", err)
		w.WriteHeader(500)
		return
	}

	finalResponse := JobListResponse{
		Status: "ok",
		Jobs:   responses,
		Count:  len(responses),
	}

	writeErr := WriteJsonContent(finalResponse, w, 200)
	if writeErr != nil {
		log.Printf("Could not write content for job list: %s", writeErr)
	}
}
