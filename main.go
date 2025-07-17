package main

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

var roomClients = make(map[string][]*websocket.Conn)
var roomClientsMutex sync.Mutex

func wsUpgrader(respWriter http.ResponseWriter, req *http.Request) (*websocket.Conn, error) {
	// Upgrade incoming to websocket
	var upgrader = websocket.Upgrader{ /*could set CrossOrigin here*/ }
	return upgrader.Upgrade(respWriter, req, nil)
}

func roomHandler(respWriter http.ResponseWriter, req *http.Request) {
	conn, upgradeErr := wsUpgrader(respWriter, req)

	if upgradeErr != nil {
		log.Println("Error upgrading connection to websocket:", upgradeErr)
		return
	}

	room := req.PathValue("roomId")

	go handleRoomConnection(conn, room)
}

func removeFromRoom(roomId string, conn *websocket.Conn) error {
	roomClientsMutex.Lock()
	defer roomClientsMutex.Unlock()

	for i, client := range roomClients[roomId] {
		if client == conn {
			roomClients[roomId] = slices.Delete(roomClients[roomId], i, i+1)
			return nil
		}
	}
	return fmt.Errorf("conn '%+v' not found in roomClients[%s]: %v", conn, roomId, roomClients[roomId])
}

func handleRoomConnection(conn *websocket.Conn, roomId string) {
	defer func() {
		_ = removeFromRoom(roomId, conn)
		conn.Close()
	}()

	roomClientsMutex.Lock()
	roomClients[roomId] = append(roomClients[roomId], conn)
	roomClientsMutex.Unlock()

	for { // Listen for incoming messages (inf. loop)
		var message Message
		readingErr := conn.ReadJSON(&message)

		if readingErr != nil {
			log.Println("Error reading message:", readingErr)
			break
		}

		log.Printf("[Room %s] Recieved message: %+v", roomId, message)

		response := message
		// response.From = "Room " + roomId
		// response.Message = fmt.Sprintf("'%s' said: %s", message.From, message.Message)

		roomClientsMutex.Lock()
		for _, roomParticipant := range roomClients[roomId] {
			writingErr := roomParticipant.WriteJSON(response)
			if writingErr != nil {
				log.Printf("Error writing message for conn '%v' in room '%s':%v", roomParticipant, roomId, writingErr)
				continue
			}
		}
		roomClientsMutex.Unlock()
	}
}

func main() {
	http.HandleFunc("/chatroom/{roomId}", roomHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "public/index.html") })
	http.HandleFunc("/terminal.css", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "public/terminal.css") })

	log.Println("Starting server on port 8080...")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Println("Error starting http server:", err)
	}
}
