package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	handler "my-go-app/internal/handler/todo"
	"my-go-app/internal/infrastructure/postgres"
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
	Update(http.ResponseWriter, *http.Request, string)
	Delete(http.ResponseWriter, *http.Request, string)
}

func main() {
	db, err := newDB()
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer db.Close()

	server := newServer(addr, db)

	go listen(server)

	waitForShutdown(server)
}

func newDB() (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func newServer(addr string, db *sql.DB) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           newMux(db),
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func newMux(db *sql.DB) http.Handler {
	repo := postgres.NewTodoPostgres(db)
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

		switch r.Method {

		case http.MethodGet:
			todoHandler.Show(w, r, id)

		case http.MethodPut:
			todoHandler.Update(w, r, id)

		case http.MethodDelete:
			todoHandler.Delete(w, r, id)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
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
