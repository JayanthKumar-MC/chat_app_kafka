package main

import (
	"chat-app/handlers"
	"chat-app/kafka"
	"chat-app/models"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	models.InitDB()
	defer models.CloseDB()

	kafka.InitProducer()
	//defer kafka.CloseProducer()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/register", handlers.ShowRegisterPage)
	r.POST("/register", handlers.Register)
	r.GET("/login", handlers.ShowLoginPage)
	r.POST("/login", handlers.Login)
	r.GET("/chat", handlers.ShowChatPage)
	r.GET("/users", handlers.ListUsers)
	r.GET("/messages", handlers.ListMessages)
	r.GET("/ws", handlers.HandleWebSocket)

	go kafka.ConsumeMessages()

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}
