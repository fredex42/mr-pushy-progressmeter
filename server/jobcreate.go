package main

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"io/ioutil"
	"log"
	"net/http"
)

type NewJobRequest struct {
	Timestamp    string `json:"timestamp"`
	UploadsCount int    `json:"uploadsCount"`
}

type NewJobReponse struct {
	Status string `json:"status"`
	JobId  string `json:"jobId"`
}

type JobCreateHandler struct {
	redisClient *redis.Client
}

func (h JobCreateHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	bodyContent, readErr := ioutil.ReadAll(request.Body)

	if AssertHttpMethod(request, w, "POST") == false {
		return
	}

	if readErr != nil {
		log.Printf("Could not read data: %s", readErr)
		w.WriteHeader(500)
		return
	}

	var req NewJobRequest
	marshalErr := json.Unmarshal(bodyContent, &req)

	if marshalErr != nil {
		log.Printf("Could not understand request data: %s", marshalErr)
		w.WriteHeader(400)
		return
	}

	uuid, createErr := CreateNewJob(h.redisClient, req)
	if createErr != nil {
		log.Printf("Could not create new job: %s", createErr)
		w.WriteHeader(500)
		return
	}

	textBytes, marshalErr := uuid.MarshalText()

	response := NewJobReponse{
		Status: "ok",
		JobId:  string(textBytes),
	}

	writeErr := WriteJsonContent(&response, w, 200)

	if writeErr != nil {
		log.Printf("Could not write response: %s", writeErr)
	}
}
