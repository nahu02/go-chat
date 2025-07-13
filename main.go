package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func wsHandler(respWriter http.ResponseWriter, req *http.Request) {
	// Upgrade incoming to websocket
	var upgrader = websocket.Upgrader{ /*could set CrossOrigin here*/ }
	conn, err := upgrader.Upgrade(respWriter, req, nil)

	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	defer conn.Close() // When the function returns, close the connection

	for { // Listen for incoming messages (inf. loop)
		_, message, readingErr := conn.ReadMessage()

		if readingErr != nil {
			log.Println("Error reading message:", readingErr)
			break
		}

		log.Println("Recieved message:", string(message))

		response := message
		writingErr := conn.WriteMessage(websocket.TextMessage, response)
		if writingErr != nil {
			log.Println("Error writing message:", writingErr)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", wsHandler)

	log.Println("Starting server on port 8000...")
	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		log.Println("Error starting http server:", err)
	}
}
