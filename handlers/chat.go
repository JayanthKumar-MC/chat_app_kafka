package handlers

import (
	"chat-app/kafka"
	"chat-app/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan models.Message)    // broadcast channel

func ShowChatPage(c *gin.Context) {
	c.HTML(http.StatusOK, "chat.html", nil)
}

func ListUsers(c *gin.Context) {
	users, err := models.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func ListMessages(c *gin.Context) {
	senderID := c.Query("sender_id")
	receiverID := c.Query("receiver_id")
	senderIDInt, err := strconv.Atoi(senderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sender ID"})
		return
	}

	receiverIDInt, err := strconv.Atoi(receiverID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receiver ID"})
		return
	}

	messages, err := models.GetMessagesBetweenUsers(senderIDInt, receiverIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

func HandleWebSocket(c *gin.Context) {
	log.Printf("Publishing message")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	clients[conn] = true

	go handleMessages()

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			delete(clients, conn)
			return
		}

		var msg models.Message
		if err := json.Unmarshal(p, &msg); err != nil {
			log.Printf("Inside Unmarshal : %v", p)
			conn.WriteMessage(websocket.TextMessage, []byte("Invalid message format"))
			continue
		}
		msg.Status = "sent"
		log.Printf("Unmarshed message : %v", msg)

		messageID, err := models.SaveMessage(msg.SenderID, msg.ReceiverID, msg.MessageText)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Could not save message"))
			continue
		}
		msg.ID = messageID

		msgBytes, err := json.Marshal(msg)
		kafka.ProduceMessage(msgBytes)
		if err != nil {
			log.Printf("Error marshaling message: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Invalid message format"))
			continue
		}

		//conn.WriteMessage(websocket.TextMessage, msgBytes)
		broadcast <- msg
	}
}

// Handle incoming messages from the broadcast channel
func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast

		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("websocket error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
