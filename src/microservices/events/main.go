package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/gorilla/mux"
)

// Event models
type UserEvent struct {
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Action    string `json:"action"`
	Timestamp string `json:"timestamp"`
}

type PaymentEvent struct {
	PaymentID  int     `json:"payment_id"`
	UserID     int     `json:"user_id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
	Timestamp  string  `json:"timestamp"`
	MethodType string  `json:"method_type"`
}

type MovieEvent struct {
	MovieID int    `json:"movie_id"`
	Title   string `json:"title"`
	Action  string `json:"action"`
	UserID  int    `json:"user_id"`
}

// EventsService handles Kafka operations
type EventsService struct {
	producer sarama.SyncProducer
	consumer sarama.Consumer
	mu       sync.RWMutex
}

var eventsService *EventsService

func main() {
	// Initialize Kafka
	if err := initKafka(); err != nil {
		log.Fatalf("Failed to initialize Kafka: %v", err)
	}
	defer eventsService.producer.Close()
	defer eventsService.consumer.Close()

	// Start consumers in background
	go startConsumers()

	// Set up HTTP routes
	router := mux.NewRouter()
	
	// Health check
	router.HandleFunc("/api/events/health", handleHealth).Methods("GET")
	
	// Event endpoints
	router.HandleFunc("/api/events/user", handleUserEvent).Methods("POST")
	router.HandleFunc("/api/events/payment", handlePaymentEvent).Methods("POST")
	router.HandleFunc("/api/events/movie", handleMovieEvent).Methods("POST")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Starting events microservice on port %s", port)
	
	// Graceful shutdown
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func initKafka() error {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		brokers = "localhost:9092"
	}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Consumer.Return.Errors = true

	// Create producer
	producer, err := sarama.NewSyncProducer(strings.Split(brokers, ","), config)
	if err != nil {
		return fmt.Errorf("failed to create producer: %w", err)
	}

	// Create consumer
	consumer, err := sarama.NewConsumer(strings.Split(brokers, ","), config)
	if err != nil {
		producer.Close()
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	eventsService = &EventsService{
		producer: producer,
		consumer: consumer,
	}

	return nil
}

func startConsumers() {
	topics := []string{"user-events", "payment-events", "movie-events"}
	
	for _, topic := range topics {
		go consumeFromTopic(topic)
	}
}

func consumeFromTopic(topic string) {
	partitionConsumer, err := eventsService.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Failed to start consumer for topic %s: %v", topic, err)
		return
	}
	defer partitionConsumer.Close()

	log.Printf("Started consumer for topic: %s", topic)

	for {
		select {
		case message := <-partitionConsumer.Messages():
			log.Printf("Received message from topic %s: %s", topic, string(message.Value))
			
		case err := <-partitionConsumer.Errors():
			log.Printf("Consumer error for topic %s: %v", topic, err)
		}
	}
}

func publishEvent(topic string, event interface{}) error {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	message := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(eventJSON),
	}

	partition, offset, err := eventsService.producer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	log.Printf("Message sent to topic %s, partition %d, offset %d", topic, partition, offset)
	return nil
}

// HTTP Handlers
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"status": true, "service": "events"})
}

func handleUserEvent(w http.ResponseWriter, r *http.Request) {
	var event UserEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := publishEvent("user-events", event); err != nil {
		log.Printf("Failed to publish user event: %v", err)
		http.Error(w, "Failed to publish event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "topic": "user-events"})
}

func handlePaymentEvent(w http.ResponseWriter, r *http.Request) {
	var event PaymentEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := publishEvent("payment-events", event); err != nil {
		log.Printf("Failed to publish payment event: %v", err)
		http.Error(w, "Failed to publish event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "topic": "payment-events"})
}

func handleMovieEvent(w http.ResponseWriter, r *http.Request) {
	var event MovieEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := publishEvent("movie-events", event); err != nil {
		log.Printf("Failed to publish movie event: %v", err)
		http.Error(w, "Failed to publish event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "topic": "movie-events"})
}
