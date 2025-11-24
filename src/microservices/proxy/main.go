package main

import (
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// Тип конфигурации для хранения переменных окружения
type Config struct {
	Port                   string
	MonolithURL           string
	MoviesServiceURL      string
	EventsServiceURL      string
	GradualMigration      bool
	MoviesMigrationPercent int
}

// ProxyService обрабатывает логику маршрутизации
type ProxyService struct {
	config           Config
	monolithProxy    *httputil.ReverseProxy
	moviesProxy      *httputil.ReverseProxy
	eventsProxy      *httputil.ReverseProxy
}

func main() {
	config := loadConfig()
	proxy := NewProxyService(config)
	
	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/health", proxy.healthHandler).Methods("GET")

	// API маршруты
	router.PathPrefix("/api/movies").HandlerFunc(proxy.moviesHandler)
	router.PathPrefix("/api/users").HandlerFunc(proxy.monolithHandler)
	router.PathPrefix("/api/payments").HandlerFunc(proxy.monolithHandler)
	router.PathPrefix("/api/subscriptions").HandlerFunc(proxy.monolithHandler)
	router.PathPrefix("/api/events").HandlerFunc(proxy.eventsHandler)

	// Обработка всех остальных маршрутов - направляем в монолит
	router.PathPrefix("/").HandlerFunc(proxy.monolithHandler)
	
	log.Printf("Starting Strangler Fig Proxy on port %s", config.Port)
	log.Printf("Gradual Migration: %v, Movies Migration Percent: %d%%", 
		config.GradualMigration, config.MoviesMigrationPercent)
	log.Fatal(http.ListenAndServe(":"+config.Port, router))
}

func loadConfig() Config {
	port := getEnv("PORT", "8000")
	monolithURL := getEnv("MONOLITH_URL", "http://localhost:8080")
	moviesServiceURL := getEnv("MOVIES_SERVICE_URL", "http://localhost:8081")
	eventsServiceURL := getEnv("EVENTS_SERVICE_URL", "http://localhost:8082")
	
	gradualMigration := getEnv("GRADUAL_MIGRATION", "false") == "true"
	moviesMigrationPercent, err := strconv.Atoi(getEnv("MOVIES_MIGRATION_PERCENT", "0"))
	if err != nil {
		moviesMigrationPercent = 0
	}
	
	return Config{
		Port:                   port,
		MonolithURL:           monolithURL,
		MoviesServiceURL:      moviesServiceURL,
		EventsServiceURL:      eventsServiceURL,
		GradualMigration:      gradualMigration,
		MoviesMigrationPercent: moviesMigrationPercent,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func NewProxyService(config Config) *ProxyService {
	monolithURL, _ := url.Parse(config.MonolithURL)
	moviesURL, _ := url.Parse(config.MoviesServiceURL)
	eventsURL, _ := url.Parse(config.EventsServiceURL)
	
	return &ProxyService{
		config:        config,
		monolithProxy: httputil.NewSingleHostReverseProxy(monolithURL),
		moviesProxy:   httputil.NewSingleHostReverseProxy(moviesURL),
		eventsProxy:   httputil.NewSingleHostReverseProxy(eventsURL),
	}
}

func (p *ProxyService) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Strangler Fig Proxy is healthy"))
}

func (p *ProxyService) moviesHandler(w http.ResponseWriter, r *http.Request) {
	// Если постепенная миграция отключена, направляем только в монолит
	if !p.config.GradualMigration {
		log.Printf("Routing movies request to monolith (migration disabled)")
		p.monolithProxy.ServeHTTP(w, r)
		return
	}

    // Упрощенный round robin
    if rand.Intn(100) < p.config.MoviesMigrationPercent {
      log.Printf("Routing movies to movies")
	  p.moviesProxy.ServeHTTP(w, r)
    } else {
      log.Printf("Routing movies to monolith")
	  p.monolithProxy.ServeHTTP(w, r)
    }
}

func (p *ProxyService) monolithHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Routing request to monolith: %s", r.URL.Path)
	p.monolithProxy.ServeHTTP(w, r)
}

func (p *ProxyService) eventsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Routing request to events service: %s", r.URL.Path)
	p.eventsProxy.ServeHTTP(w, r)
}
