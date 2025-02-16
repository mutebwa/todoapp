package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mutebwa/todoapp/handlers"
)

func main() {
	// Configure server settings from environment variables
	port := getEnv("PORT", "8088")
	shutdownTimeout := getEnvDuration("SHUTDOWN_TIMEOUT", 5*time.Second)
	readTimeout := getEnvDuration("READ_TIMEOUT", 5*time.Second)
	writeTimeout := getEnvDuration("WRITE_TIMEOUT", 10*time.Second)
	idleTimeout := getEnvDuration("IDLE_TIMEOUT", 15*time.Second)

	// Create router with middleware chain
	router := http.NewServeMux()
	router.HandleFunc("/add", handlers.AddTask)
	router.HandleFunc("/tasks", handlers.GetTasks)

	// Static file handling with SPA support
	spa := spaHandler{staticPath: "public", indexPath: "index.html"}
	router.Handle("/", securityHeaders(spa))

	// Create server with configured timeouts
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           loggingMiddleware(router),
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	// Graceful shutdown channel
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig
		log.Println("Initiating graceful shutdown...")

		shutdownCtx, cancel := context.WithTimeout(serverCtx, shutdownTimeout)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("Graceful shutdown timed out.. forcing exit")
			}
		}()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("Server shutdown failed: %v", err)
		}
		serverStopCtx()
	}()

	// Start server
	log.Printf("Starting server on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	log.Println("Server stopped")
}

// Helper functions and handlers

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		duration, err := time.ParseDuration(value)
		if err == nil {
			return duration
		}
	}
	return defaultValue
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}

// spaHandler implements the http.Handler interface for SPA support
type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check if the requested file exists
	path := h.staticPath + r.URL.Path
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// File does not exist, serve index.html
		http.ServeFile(w, r, h.staticPath+"/"+h.indexPath)
		return
	} else if err != nil {
		// Other error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Otherwise, use http.FileServer to serve static files
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}