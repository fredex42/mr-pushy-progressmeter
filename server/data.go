package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"log"
)

type JobEntry struct {
	JobId        uuid.UUID `json:"jobId"`
	Timestamp    string    `json:"timestamp"`
	UploadsCount int       `json:"uploadsCount"`
}

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

	entry := JobEntry{
		JobId:        jobId,
		Timestamp:    request.Timestamp,
		UploadsCount: request.UploadsCount,
	}

	jsonContent, _ := json.Marshal(&entry)

	client.Set(fmt.Sprintf("job:%s", jobId), jsonContent, -1)

	return jobId, nil
}

func retrieveKey(client *redis.Client, jobKey string) (*JobEntry, error) {
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

	var entry JobEntry
	unmarshalErr := json.Unmarshal(bytesContent, &entry)

	if unmarshalErr != nil {
		log.Printf("Could not unmarshal content from database: %s", unmarshalErr)
		return nil, unmarshalErr
	}

	return &entry, nil
}

func GetJob(client *redis.Client, jobId string) (*JobEntry, error) {
	jobKey := fmt.Sprintf("job:%s", jobId)
	return retrieveKey(client, jobKey)
}

func ListJobs(client *redis.Client, limit int64) ([]*JobEntry, error) {
	var cursor uint64
	keys, _, err := client.Scan(cursor, "job:*", limit).Result()

	rtn := make([]*JobEntry, len(keys))

	if err != nil {
		log.Printf("Could not scan database for jobs: %s", err)
		return rtn, err
	}

	for ctr, keyString := range keys {
		entryContent, retrieveErr := retrieveKey(client, keyString)
		if retrieveErr != nil {
			log.Printf("Could not retrieve item %d of %d with key %s: %s", ctr, len(keys), keyString, retrieveErr)
			return rtn, err
		}
		rtn[ctr] = entryContent
	}
	return rtn, nil
}
