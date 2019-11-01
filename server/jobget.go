package main

import (
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"log"
	"net/http"
	"net/url"
)

type JobGetResponse struct {
	Status string    `json:"status"`
	Entry  *JobEntry `json:"entry"`
}

type JobGetHandler struct {
	redisClient *redis.Client
}

func (h JobGetHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	requestUri, uriParseErr := url.ParseRequestURI(request.RequestURI)

	if uriParseErr != nil {
		log.Printf("Could not understand incoming request URI '%s': %s", request.RequestURI, uriParseErr)
		response := GenericErrorResponse{
			Status: "error",
			Detail: "invalid URI",
		}
		WriteJsonContent(&response, w, 400)
		return
	}

	queryParams := requestUri.Query()
	jobIdString := queryParams.Get("jobId")

	jobId, uuidParseErr := uuid.Parse(jobIdString)
	if uuidParseErr != nil {
		log.Printf("Could not parse job ID string '%s' into a UUID: %s", jobIdString, uuidParseErr)
		response := GenericErrorResponse{
			Status: "error",
			Detail: "malformed UUID",
		}
		WriteJsonContent(&response, w, 400)
		return
	}

	jobInfo, getErr := GetJob(h.redisClient, jobId.String())

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
