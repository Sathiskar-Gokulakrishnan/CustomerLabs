package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var WEBHOOK_URL string = "https://webhook.site/fc00fb8b-98bf-4992-a55d-5cbbe7d00d6f"

type EventData struct {
	Ev             string                    `json:"event"`
	Et             string                    `json:"event_type"`
	ID             string                    `json:"app_id"`
	UID            string                    `json:"user_id"`
	MID            string                    `json:"message_id"`
	T              string                    `json:"page_title"`
	P              string                    `json:"page_url"`
	L              string                    `json:"browser_language"`
	SC             string                    `json:"screen_size"`
	Attributes     map[string]Attributes     `json:"attributes"`
	UserAttributes map[string]UserAttributes `json:"traits"`
}

type Attributes struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type UserAttributes struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

func (s *EventData) UnmarshalJSON(data []byte) error {
	var aux map[string]interface{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	s.Attributes = make(map[string]Attributes)
	s.UserAttributes = make(map[string]UserAttributes)

	var attributes Attributes
	var userAttributes UserAttributes

	// Extract values
	for key, value := range aux {
		if key == "ev" {
			s.Ev = value.(string)
		} else if key == "et" {
			s.Et = value.(string)
		} else if key == "id" {
			s.ID = value.(string)
		} else if key == "uid" {
			s.UID = value.(string)
		} else if key == "mid" {
			s.MID = value.(string)
		} else if key == "t" {
			s.T = value.(string)
		} else if key == "p" {
			s.P = value.(string)
		} else if key == "l" {
			s.L = value.(string)
		} else if key == "sc" {
			s.SC = value.(string)
		} else if strings.HasPrefix(key, "atrk") {
			id := strings.Split(key, "k")[1]
			attributes.Value = aux["atrv"+id].(string)
			attributes.Type = aux["atrt"+id].(string)
			s.Attributes[value.(string)] = attributes
		} else if strings.HasPrefix(key, "uatrk") {
			id := strings.Split(key, "k")[1]
			userAttributes.Value = aux["uatrv"+id].(string)
			userAttributes.Type = aux["uatrt"+id].(string)
			s.UserAttributes[value.(string)] = userAttributes
		}
	}

	return nil
}

func worker(chanReq chan *http.Request, chanResp chan []byte) {
	fmt.Println("Worker Starts")
	reqData := <-chanReq
	var eventData EventData

	decoder := json.NewDecoder(reqData.Body)

	err := decoder.Decode(&eventData)
	if err != nil {
		panic(err)
	}

	data, _ := json.Marshal(eventData)
	chanResp <- data
}

func requestHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		fmt.Fprint(w, "METHOD NOT ALLOWED")
		return
	}

	chanReq := make(chan *http.Request)
	chanResp := make(chan []byte)

	go worker(chanReq, chanResp)

	chanReq <- req
	getFromChannel := <-chanResp
	reader := bytes.NewReader(getFromChannel)

	http.Post(WEBHOOK_URL, "text/json", reader)
	fmt.Fprint(w, "Success")
	return
}

func main() {
	fmt.Println("Hello World")
	http.HandleFunc("/", requestHandler)
	http.ListenAndServe(":80", nil)
}

/*
TODO

1. Create a HTTP server in golang, that will receive request in below format [DONE]
2. Create a Golang channel to send this request to a golang worker [DONE]
3. Create a go worker that will receive the above message from the channel and convert into
below format[Done]
*/
