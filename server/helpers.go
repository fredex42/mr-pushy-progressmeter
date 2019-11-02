package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func WriteJsonContent(content interface{}, w http.ResponseWriter, statusCode int) error {
	contentBytes, marshalErr := json.Marshal(content)
	if marshalErr != nil {
		return marshalErr
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Content-Length", strconv.FormatInt(int64(len(contentBytes)), 10))
	w.WriteHeader(statusCode)
	_, writeErr := w.Write(contentBytes)
	return writeErr
}

func AssertHttpMethod(request *http.Request, w http.ResponseWriter, method string) bool {
	if request.Method != method {
		log.Printf("Got a %s request, expecting %s", request.Method, method)
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(405)
		w.Write([]byte(fmt.Sprintf("Invalid method, expecting %s", method)))
		return false
	} else {
		return true
	}
}

/**
Breaks down the incoming request URI into a map of string->string
*/
func GetQueryParams(incomingRequestUri string) (*url.Values, error) {
	requestUri, uriParseErr := url.ParseRequestURI(incomingRequestUri)

	if uriParseErr != nil {
		log.Printf("Could not understand incoming request URI '%s': %s", incomingRequestUri, uriParseErr)
		return nil, errors.New("Invalid URI")
	}

	rtn := requestUri.Query()
	return &rtn, nil
}

/**
gets just the "JobID" parameter from the provided query string and returns it as a pointer to UUID
if it does not exist or is not a valid UUID, a GenericErrorResponse object is returned that is suitable
to be written directly to the outgoing response.

This is a convenience function that calls GetQueryParams and GetJobIdFromValues
*/
func GetJobIdFromQuerystring(incomingRequestUri string) (*uuid.UUID, *GenericErrorResponse) {
	queryParams, err := GetQueryParams(incomingRequestUri)
	if err != nil {
		return nil, &GenericErrorResponse{
			Status: "error",
			Detail: err.Error(),
		}
	}
	return GetJobIdFromValues(queryParams)
}

func GetJobIdFromValues(queryParams *url.Values) (*uuid.UUID, *GenericErrorResponse) {
	jobIdString := queryParams.Get("jobId")

	jobId, uuidParseErr := uuid.Parse(jobIdString)
	if uuidParseErr != nil {
		log.Printf("Could not parse job ID string '%s' into a UUID: %s", jobIdString, uuidParseErr)
		return nil, &GenericErrorResponse{
			Status: "error",
			Detail: "malformed UUID",
		}
	}
	return &jobId, nil
}
