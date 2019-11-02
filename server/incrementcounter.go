package main

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"log"
	"net/http"
)

//can't declare this as a constant, unfortunately
var counterNames = []string{"UploadSuccess", "UploadFail", "VerifySuccessDelete", "VerifySuccessKeep", "VerifyFailure"}

type IncrementCounterHandler struct {
	redisClient *redis.Client
}

func IsValidCountername(toTest string) bool {
	for _, name := range counterNames {
		if toTest == name {
			return true
		}
	}
	return false
}

func (h IncrementCounterHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	queryParams, qpErr := GetQueryParams(req.RequestURI)

	if qpErr != nil {
		log.Printf("Could not get query params: %s", qpErr)
		response := &GenericErrorResponse{
			Status: "error",
			Detail: "Could not get query params",
		}
		WriteJsonContent(response, w, 400)
		return
	}

	jobId, uuidErr := GetJobIdFromValues(queryParams)

	if uuidErr != nil {
		WriteJsonContent(uuidErr, w, 400)
		return
	}

	counterName := queryParams.Get("counter")

	if IsValidCountername(counterName) == false {
		log.Printf("%s is not a valid counter name. Must be one of %s", counterName, counterNames)
		response := &InvalidOptionResponse{
			Status:  "error",
			Detail:  fmt.Sprintf("%s is not a valid counter name.", counterName),
			Options: counterNames,
		}
		WriteJsonContent(response, w, 400)
		return
	}

	_, err := IncrementCounter(h.redisClient, *jobId, counterName)

	if err != nil {
		log.Printf("Could not increment counter '%s': %s", counterName, err)
		response := &GenericErrorResponse{
			Status: "error",
			Detail: "Could not increment database counter",
		}
		WriteJsonContent(response, w, 500)
		return
	}

	response := &GenericErrorResponse{
		Status: "ok",
		Detail: "incremented counter",
	}
	WriteJsonContent(response, w, 200)
	return
}
