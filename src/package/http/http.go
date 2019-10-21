package client

import (
	"log"
	"net/http"
)

type Request struct {
	*http.Request
}

func Fire(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
		return nil, err
	}

	return resp, nil
}
