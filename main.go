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
	// Create a new ServeMux.
	mux := http.NewServeMux()
	
	// API endpoints.
	mux.HandleFunc("/add", handlers.AddTask)
	mux.HandleFunc("/tasks", handlers.GetTasks)
	
	// Register a handler to serve the HTML file.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Serve the "index.html" file located in the "views" folder.
		http.ServeFile(w, r, "public/index.html")
	})

	// Create a custom HTTP server with timeouts.
	server := &http.Server{
		Addr:         ":8088",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Start the server in a separate goroutine.
	go func() {
		log.Println("To-Do App starting on :8088")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP Server error: %v", err)
		}
	}()

	// Listen for OS signals to gracefully shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for in-flight requests.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %v", err)
	}
	log.Println("Server gracefully stopped")
}