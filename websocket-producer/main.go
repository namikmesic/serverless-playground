package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	url := "wss://7dsv3kb2al.execute-api.us-west-2.amazonaws.com/Prod"

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket:", err)
	}
	defer conn.Close()

	log.Println("Connected to WebSocket")

	for {
		message := struct {
			Action string `json:"action"`
			Data   struct {
				Message   string `json:"message"`
				Timestamp int64  `json:"timestamp"`
			} `json:"data"`
		}{
			Action: "sendmessage",
			Data: struct {
				Message   string `json:"message"`
				Timestamp int64  `json:"timestamp"`
			}{
				Message:   "hello world",
				Timestamp: time.Now().UnixNano(),
			},
		}

		jsonMessage, err := json.Marshal(message)
		if err != nil {
			log.Println("Error encoding message to JSON:", err)
			continue
		}

		err = conn.WriteMessage(websocket.TextMessage, jsonMessage)
		if err != nil {
			log.Println("Error sending message to WebSocket:", err)
			continue
		}

		log.Println("Sent message to WebSocket:", string(jsonMessage))

		time.Sleep(200 * time.Millisecond)
	}
}
