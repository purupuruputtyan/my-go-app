package main

import (
	"context"
	"log"
	"my-go-app/internal/handler/todo"
	"my-go-app/internal/infrastructure/todo"
	"my-go-app/internal/usecase/todo"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	repo := memory.NewTodoMemory()
	todoUsecase := todo.NewTodoUseCase(repo)
	todoHandler := handler.New(todoUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			todoHandler.FindAll(w, r)

		case http.MethodPost:
			todoHandler.Create(w, r)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"Go API is running"}`))
	})

	addr := ":8080"
	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("API server started on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Println("shutting down server...")

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}

	log.Println("server stopped")
}
