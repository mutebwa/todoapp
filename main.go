package main

import (
	// "fmt"
	"log"
	"net"
	"net/http"
)

func main() {
	// Start HTTP server for serving HTML
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "views/index.html")
		})
		log.Println("To-Do App starting on :8088")
		err := http.ListenAndServe(":8088", nil)
		if err != nil {
			log.Fatal("HTTP Server error: ", err)
		}
	}()

	// Start TCP server
	go func() {
		listener, err := net.Listen("tcp", ":1234")
		if err != nil {
			log.Fatal("TCP Server error: ", err)
		}
		defer listener.Close()
		log.Println("TCP server starting on :1234")

		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("TCP connection error: ", err)
				continue
			}
			go handleConnection(conn)
		}
	}()

	// Keep main goroutine running
	select {}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	conn.Write([]byte("Hello from TCP server!\n"))
}