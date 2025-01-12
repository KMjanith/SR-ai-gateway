package main

import (
	"api-gateway/handlers"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	godotenv.Load() //make the connection to the mongodb
	log.Println("starting Api gateway")

	ctx, cancel := context.WithTimeout(context.Background(), 10000000*time.Second) // Set up a context with timeout
	defer cancel()

	//---------------------------------------RABBITMQ CONNECTION---------------------------------//

	rabbitUrl := os.Getenv("RABBITMQ_URI")
	conn, err := amqp.Dial(rabbitUrl) //making the connection to the rabbitmq
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ctx = context.WithValue(ctx, "conn", conn) //adding the connection to the context

	consumerChannel, err := conn.Channel() //Open a channel for consumer
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer consumerChannel.Close()

	ctx = context.WithValue(ctx, "consumerChannel", consumerChannel) //adding the consumer channel to the context

	producerChannel, err := conn.Channel() //open a channel for producer
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer producerChannel.Close()

	log.Println("Connected to RabbitMQ")

	ctx = context.WithValue(ctx, "producerChannel", producerChannel) //adding the producer channel to the context

	// Setup routes
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterHandler(w, r, ctx)
	})

	//sorting
	// Setup routes
	http.HandleFunc("/sort", func(w http.ResponseWriter, r *http.Request) {
		handlers.SortingHandler(w, r, ctx)
	})

	// Start the server
	port := ":8081"
	fmt.Printf("Server running on http://api-gateway%s\n", port)
	log.Fatal(http.ListenAndServe(port, http.DefaultServeMux))
}
