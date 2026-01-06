package main

import (
	"log"
	"net/http"
	"os"
	"simple-api/internal/database"
	"simple-api/internal/handlers"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://taskuser:taskpass@localhost:5432/tasksdb?sslmode=disable"
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	log.Printf("Connecting to server at port %s", serverPort)
	db, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	taskscore := database.NewTaskScore(db)
	log.Printf("connected to db")

	handler := handlers.NewHandler(taskscore)
	mux := http.NewServeMux()

	mux.HandleFunc("/tasks/create", methodHandler(handler.CreateTask, http.MethodPost))

	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/tasks/" {
			if r.Method == http.MethodGet {
				handler.GetAllTasks(w, r)
			} else {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			}
			return
		}
		taskIDHandler(handler)(w, r)
	})

	loggedMux := loggingMiddleware(mux)
	serverAddr := ":" + serverPort
	log.Printf("Server starting on %s", serverAddr)
	err = http.ListenAndServe(serverAddr, loggedMux)
	if err != nil {
		log.Fatal(err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func methodHandler(handler http.HandlerFunc, allowedMethod string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}

func taskIDHandler(handler *handlers.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetTaskByID(w, r)
		case http.MethodPut:
			handler.UpdateTask(w, r)
		case http.MethodDelete:
			handler.DeleteTask(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}
