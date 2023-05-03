package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	url := "wss://7dsv3kb2al.execute-api.us-west-2.amazonaws.com/Prod"

	for {
		log.Println("Connecting to", url)
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			log.Println("Error connecting to WebSocket:", err)
			log.Println("Retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}
		defer conn.Close()

		log.Println("Connected to WebSocket")

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading message from WebSocket:", err)
				break
			}

			log.Printf("Received message: %s\n", message)
		}

		log.Println("Connection to WebSocket closed")
		log.Println("Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}
