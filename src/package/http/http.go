package client

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Request struct {
	*http.Request
}

func Fire(req *http.Request) (*http.Response, error) {
	client := &http.Client{}

	log.Println("Client Fired")

	return retry(5, time.Second, func() (*http.Response, error) {
		resp, err := client.Do(req)

		if err != nil {
			log.Println("Error on response.\n[ERROR] -", err)
			return resp, err
		}

		//defer resp.Body.Close()

		switch {
		case strings.Contains(resp.Status, "404"):
			return resp, fmt.Errorf("resource not found: %v", resp.Status)
		default:
			return resp, nil
		}
	})
}

func retry(attempts int, sleep time.Duration, fn func() (*http.Response, error)) (*http.Response, error) {
	response, err := fn()

	if err != nil {
		if s, ok := err.(stop); ok {
			return response, s.error
		}

		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			log.Println("HTTP Request failed - Retry attempt " + string(attempts))
			return retry(attempts, 2*sleep, fn)
		}
		return response, err
	}

	return response, nil
}

type stop struct {
	error
}

type Response struct {
	Errors []Errors `json:"errors"`
}

type Errors struct {
	Message string `json:"message"`
	Code string `json:"code"`
	HttpStatus int `json:"http_status"`
}