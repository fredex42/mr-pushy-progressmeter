package main

import (
	"github.com/go-redis/redis/v7"
	"log"
	"net/http"
)

type JobGetResponse struct {
	Status string    `json:"status"`
	Entry  *JobEntry `json:"entry"`
}

type JobGetHandler struct {
	redisClient *redis.Client
}

func (h JobGetHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	jobId, errorResponse := GetJobIdFromQuerystring(request.RequestURI)

	if errorResponse != nil {
		WriteJsonContent(errorResponse, w, 400)
		return
	}

	jobInfo, getErr := GetJob(h.redisClient, *jobId)

	if getErr != nil {
		response := GenericErrorResponse{
			Status: "error",
			Detail: "Could not get job data from db",
		}
		WriteJsonContent(&response, w, 500)
		return
	}

	response := JobGetResponse{
		Status: "ok",
		Entry:  jobInfo,
	}
	writeErr := WriteJsonContent(&response, w, 200)
	if writeErr != nil {
		log.Printf("Could not write response: %s", writeErr)
	}

}
