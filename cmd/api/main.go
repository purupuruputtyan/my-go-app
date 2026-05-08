package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	handler "my-go-app/internal/handler/todo"
	memory "my-go-app/internal/infrastructure/todo"
	usecase "my-go-app/internal/usecase/todo"
)

const (
	addr            = ":8080"
	shutdownTimeout = 10 * time.Second
)

type todoRoutesHandler interface {
	FindAll(http.ResponseWriter, *http.Request)
	Create(http.ResponseWriter, *http.Request)
	Show(http.ResponseWriter, *http.Request, string)
}

func main() {
	server := newServer(addr)

	go listen(server)

	waitForShutdown(server)
}

func newServer(addr string) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           newMux(),
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func newMux() http.Handler {
	repo := memory.NewTodoMemory()
	todoUsecase := usecase.NewTodoUseCase(repo)
	todoHandler := handler.New(todoUsecase)

	mux := http.NewServeMux()

	registerRootRoute(mux)
	registerHealthRoutes(mux)
	registerTodoRoutes(mux, todoHandler)

	return mux
}

func registerRootRoute(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"Go API is running"}`))
	})
}

func registerHealthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
}

func registerTodoRoutes(mux *http.ServeMux, todoHandler todoRoutesHandler) {
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

	mux.HandleFunc("/todos/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/todos/")
		if id == "" {
			http.NotFound(w, r)
			return
		}

		todoHandler.Show(w, r, id)
	})
}

func listen(server *http.Server) {
	log.Printf("API server started on %s", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func waitForShutdown(server *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	log.Println("shutting down server...")

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}

	log.Println("server stopped")
}
