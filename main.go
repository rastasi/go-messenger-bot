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
	VERIFY_TOKEN      = "rubyforever"
	PAGE_ACCESS_TOKEN = "EAACChyvB7ZCMBAPXK1i9VDFJyPdEEArsMuAuGNFsuCvfe9Sp6LCX21Rkmii580x25GKUVrEYQkx7BDPdOgnibTt9I9jIx4uIgduUbaafmkmiBGn01KxhIU8yKaXjpBQfRZBDhmy2oiP1wsqh0CKgoDZCj6LEGI5mVD3BcZBtBAZDZD"
)

type Participant struct {
	Id string `json:"id"`
}

type Message struct {
	Text string `json:"text"`
}

type Response struct {
	Recipient Participant `json:"recipient"`
	Message   Message     `json:"message"`
}

type Messaging struct {
	Sender    Participant `json:"sender"`
	Recipient Participant `json:"recipient"`
	Timestamp int         `json:"timestamp"`
	Message   Message     `json:"message"`
}

type Entry struct {
	Id        string      `json:"id"`
	Time      int         `json:"time"`
	Messaging []Messaging `json:"messaging"`
}

type Answer struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

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
		if mode == "subscribe" && token == VERIFY_TOKEN {
			w.Write([]byte(challenge))
		} else {
			w.Write([]byte("Forbidden"))
		}
	}
}

func send_response(sender_id string, text string) {
	var res = Response{Recipient: Participant{Id: sender_id}, Message: Message{Text: text}}

	res_marshalled, err := json.Marshal(res)

	if err != nil {
		fmt.Println(err)
	}

	url := "https://graph.facebook.com/v2.6/me/messages"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(res_marshalled))
	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("access_token", PAGE_ACCESS_TOKEN)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func postWebhook(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var answer Answer
	err = json.Unmarshal(b, &answer)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if answer.Object == "page" {
		event := answer.Entry[0].Messaging[0]
		sender_id := event.Sender.Id
		text := event.Message.Text
		send_response(sender_id, text)
	}
	w.Write([]byte("EVENT_RECEIVED"))
}
