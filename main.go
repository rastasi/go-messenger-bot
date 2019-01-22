package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rastasi/go-messenger-bot/models"
)

const (
	verifyToken     = "goforever"
	pageAccessToken = "EAACChyvB7ZCMBAPXK1i9VDFJyPdEEArsMuAuGNFsuCvfe9Sp6LCX21Rkmii580x25GKUVrEYQkx7BDPdOgnibTt9I9jIx4uIgduUbaafmkmiBGn01KxhIU8yKaXjpBQfRZBDhmy2oiP1wsqh0CKgoDZCj6LEGI5mVD3BcZBtBAZDZD"
	messagesURL     = "https://graph.facebook.com/v2.6/me/messages"
)

var (
	address string
)

func init() {
	flag.StringVar(&address, "address", ":3000", "")
	flag.Parse()
}

func main() {
	router := chi.NewRouter()
	router.Get("/webhook", getWebhook)
	router.Post("/webhook", postWebhook)
	log.Printf("Listening on %s", address)
	log.Fatal(http.ListenAndServe(address, router))
}

func getWebhook(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")

	if mode != "" && token != "" {
		if mode == "subscribe" && token == verifyToken {
			w.Write([]byte(challenge))
		} else {
			w.Write([]byte("Forbidden"))
		}
	}
}

func getJSONResponse(response models.Response) []byte {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Print(err)
	}
	return jsonResponse
}

func sendResponse(senderID string, text string) {
	response := models.Response{Recipient: models.Participant{Id: senderID}, Message: models.Message{Text: text}}
	jsonResponse := getJSONResponse(response)

	request, err := http.NewRequest("POST", messagesURL, bytes.NewBuffer(jsonResponse))
	request.Header.Set("Content-Type", "application/json")

	request.URL.Query().Add("access_token", pageAccessToken)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Print(err)
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

func getAnswer(w http.ResponseWriter, body []byte) models.Answer {
	var answer models.Answer
	err := json.Unmarshal(body, &answer)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	return answer
}

func getResponseText(input string) string {
	switch input {
	case "hello":
		return "hello!"
	case "help":
		return "Can I help you?"
	case "contact":
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
