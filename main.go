package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

func wsHandler(respWriter http.ResponseWriter, req *http.Request) {
	// Upgrade incoming to websocket
	var upgrader = websocket.Upgrader{ /*could set CrossOrigin here*/ }
	conn, err := upgrader.Upgrade(respWriter, req, nil)

	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	go handleConnection(conn)
}

func handleConnection(conn *websocket.Conn) {
	defer conn.Close() // When the function returns, close the connection

	for { // Listen for incoming messages (inf. loop)
		var message Message
		readingErr := conn.ReadJSON(&message)

		if readingErr != nil {
			log.Println("Error reading message:", readingErr)
			break
		}

		log.Printf("Recieved message: %+v", message)

		response := message
		response.From = "Server"
		response.Message = fmt.Sprintf("'%s' said: %s", message.From, message.Message)

		writingErr := conn.WriteJSON(response)
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
