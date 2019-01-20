package main

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
