package main

import (
	"github.com/go-redis/redis/v7"
	"log"
	"net/http"
)

type JobCompleteHandler struct {
	redisClient *redis.Client
}

func (h JobCompleteHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	if AssertHttpMethod(request, w, "PUT") == false {
		return
	}

	jobId, errorResponse := GetJobIdFromQuerystring(request.RequestURI)

	if errorResponse != nil {
		WriteJsonContent(errorResponse, w, 400)
		return
	}

	setErr := SetJobStatus(h.redisClient, *jobId, CompletedSuccess)
	if setErr != nil {
		log.Printf("Could not set success for %s: %s", jobId.String(), setErr)
		response := GenericErrorResponse{
			Status: "error",
			Detail: "Could not write update to database",
		}
		WriteJsonContent(&response, w, 500)
		return
	}
}
