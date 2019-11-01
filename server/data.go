package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"log"
)

const InProgress = "InProgress"
const CompletedSuccess = "CompletedSuccess"
const CompletedFailure = "CompletedFailure"

type JobEntry struct {
	JobId        uuid.UUID `json:"jobId"`
	Timestamp    string    `json:"timestamp"`
	UploadsCount int       `json:"uploadsCount"`
	Status       string    `json:"status"`
}

type Counter struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
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
		Status:       InProgress,
	}

	jsonContent, _ := json.Marshal(&entry)

	client.Set(fmt.Sprintf("job:%s", jobId), jsonContent, -1)

	return jobId, nil
}

func SetJobStatus(client *redis.Client, jobId uuid.UUID, newStatus string) error {
	currentJob, getErr := GetJob(client, jobId)

	if getErr != nil {
		return getErr
	}

	currentJob.Status = newStatus
	jsonContent, _ := json.Marshal(&currentJob)
	return client.Set(fmt.Sprintf("job:%s", jobId.String()), jsonContent, -1).Err()
}

func retrieveJobKey(client *redis.Client, jobKey string) (*JobEntry, error) {
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

func GetJob(client *redis.Client, jobId uuid.UUID) (*JobEntry, error) {
	jobKey := fmt.Sprintf("job:%s", jobId.String())
	return retrieveJobKey(client, jobKey)
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
		entryContent, retrieveErr := retrieveJobKey(client, keyString)
		if retrieveErr != nil {
			log.Printf("Could not retrieve item %d of %d with key %s: %s", ctr, len(keys), keyString, retrieveErr)
			return rtn, err
		}
		rtn[ctr] = entryContent
	}
	return rtn, nil
}

func IncrementCounter(client *redis.Client, jobId uuid.UUID, counterName string) (int64, error) {
	keyName := fmt.Sprintf("%s:%s", jobId.String(), counterName)
	return client.Incr(keyName).Result()
}

func GetCounter(client *redis.Client, jobId uuid.UUID, counterName string) (int64, error) {
	keyName := fmt.Sprintf("%s:%s", jobId.String(), counterName)
	return client.Get(keyName).Int64()
}

func ListCounters(client *redis.Client, jobId uuid.UUID) ([]*Counter, error) {
	var cursor uint64
	keySearchTerm := fmt.Sprintf("%s:*", jobId.String())
	keys, cursor, err := client.Scan(cursor, keySearchTerm, 100).Result()

	rtn := make([]*Counter, len(keys))

	if err != nil {
		log.Printf("Could not scan database for counters: %s", err)
		return rtn, err
	}

	for ctr, keyString := range keys {
		getResult, getErr := client.Get(keyString).Int64()
		if getErr != nil {
			log.Printf("Could not get counter key %s (%d of %d): %s", keyString, ctr, len(keys), getErr)
			return rtn, getErr
		}
		rtn[ctr] = &Counter{
			Name:  keyString,
			Value: getResult,
		}
	}
	return rtn, nil
}
