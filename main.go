package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"bytes"
)

const VERIFY_TOKEN = "rubyforever"
const PAGE_ACCESS_TOKEN = "EAACChyvB7ZCMBAPXK1i9VDFJyPdEEArsMuAuGNFsuCvfe9Sp6LCX21Rkmii580x25GKUVrEYQkx7BDPdOgnibTt9I9jIx4uIgduUbaafmkmiBGn01KxhIU8yKaXjpBQfRZBDhmy2oiP1wsqh0CKgoDZCj6LEGI5mVD3BcZBtBAZDZD"

type Recipient struct {
	Id int64 `json:"id"`
}

type Message struct {
	Text string `json:"text"`
}

type Response struct {
  Recipient Recipient `json:"recipient"`
  Message Message `json:"message"`
}

func main() {
	var router = mux.NewRouter()
	router.HandleFunc("/webhook", getWebhook).Methods("GET")
	router.HandleFunc("/webhook", postWebhook).Methods("POST")
	fmt.Println("Running server!")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func getWebhook(w http.ResponseWriter, r *http.Request) {
	var mode = r.URL.Query()["hub.mode"][0]
  var token = r.URL.Query()["hub.verify_token"][0]
	var challenge = r.URL.Query()["hub.challenge"][0]
	
	if mode != "" && token != "" {
    if mode == "subscribe" && token == VERIFY_TOKEN {
      fmt.Println("WEBHOOK_VERIFIED");
			w.Write([]byte(challenge))
    } else {
      w.Write([]byte("Forbidden"));
    }
  }
}

func send_response(sender_id int64, text string) {
	var res = Response{ Recipient: Recipient{ Id: sender_id }, Message: Message{ Text: text } }

	res_marshalled, err := json.Marshal(res)
	
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res_marshalled)

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
	var sender_id = int64(1589487791148693)
	var text = string("Message answer")
	
	send_response(sender_id, text)

	w.Write([]byte("EVENT_RECEIVED"));
}