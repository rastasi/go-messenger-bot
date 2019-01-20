package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	verifyToken     = "rubyforever"
	pageAccessToken = "EAACChyvB7ZCMBAPXK1i9VDFJyPdEEArsMuAuGNFsuCvfe9Sp6LCX21Rkmii580x25GKUVrEYQkx7BDPdOgnibTt9I9jIx4uIgduUbaafmkmiBGn01KxhIU8yKaXjpBQfRZBDhmy2oiP1wsqh0CKgoDZCj6LEGI5mVD3BcZBtBAZDZD"
	messagesURL     = "https://graph.facebook.com/v2.6/me/messages"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/webhook", getWebhook).Methods("GET")
	router.HandleFunc("/webhook", postWebhook).Methods("POST")
	fmt.Println("Running server!")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func getWebhook(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query()["hub.mode"][0]
	token := r.URL.Query()["hub.verify_token"][0]
	challenge := r.URL.Query()["hub.challenge"][0]

	if mode != "" && token != "" {
		if mode == "subscribe" && token == verifyToken {
			w.Write([]byte(challenge))
		} else {
			w.Write([]byte("Forbidden"))
		}
	}
}

func getJSONResponse(response Response) []byte {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}
	return jsonResponse
}

func sendResponse(senderID string, text string) {
	response := Response{Recipient: Participant{Id: senderID}, Message: Message{Text: text}}
	jsonResponse := getJSONResponse(response)

	request, err := http.NewRequest("POST", messagesURL, bytes.NewBuffer(jsonResponse))
	request.Header.Set("Content-Type", "application/json")

	q := request.URL.Query()
	q.Add("access_token", pageAccessToken)
	request.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func getBody(w http.ResponseWriter, r *http.Request) []byte {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	return body
}

func getAnswer(w http.ResponseWriter, body []byte) Answer {
	var answer Answer
	err := json.Unmarshal(body, &answer)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	return answer
}

func getResponseText(input string) string {
	if input == "hello" {
		return "hello!"
	}
	if input == "help" {
		return "Can I help you?"
	}
	if input == "contact" {
		return "www.aiventure.me"
	}
	return "Sorry, I don't undersand."
}

func postWebhook(w http.ResponseWriter, r *http.Request) {
	body := getBody(w, r)
	answer := getAnswer(w, body)

	if answer.Object == "page" {
		event := answer.Entry[0].Messaging[0]
		if event.Message.Text != "" {
			senderID := event.Sender.Id
			text := getResponseText(event.Message.Text)
			sendResponse(senderID, text)
		}
	}
	w.Write([]byte("EVENT_RECEIVED"))
}
