package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"log"
)

func findAvailableAddress(client *redis.Client) uuid.UUID {
	u, _ := uuid.NewRandom()

	existResult := client.Exists(fmt.Sprintf("job:%s", u))

	if existResult.Val() == 1 {
		return findAvailableAddress(client)
	} else {
		return u
	}
}

func CreateNewJob(client *redis.Client, request NewJobRequest) (uuid.UUID, error) {
	jobId := findAvailableAddress(client)

	jsonContent, _ := json.Marshal(&request)

	client.Set(fmt.Sprintf("job:%s", jobId), jsonContent, -1)

	return jobId, nil
}

func GetJob(client redis.Client) (*NewJobRequest, error) {
	jobKey := fmt.Sprintf("job:%s")
	result := client.Get(jobKey)

	if result.Err() != nil {
		log.Printf("Could not get retrieve %s", jobKey)
		return nil, result.Err()
	}

	bytesContent, getErr := result.Bytes()
	if getErr != nil {
		log.Printf("Could not get bytes from result: %s", getErr)
		return nil, getErr
	}

	var originalRequet NewJobRequest
	unmarshalErr := json.Unmarshal(bytesContent, &originalRequet)

	if unmarshalErr != nil {
		log.Printf("Could not unmarshal content from database: %s", unmarshalErr)
		return nil, unmarshalErr
	}

	return &originalRequet, nil
}
