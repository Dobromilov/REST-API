package handlers

import (
	"encoding/json"
	"net/http"
	"simple-api/internal/database"
	"simple-api/internal/models"
	"strconv"
	"strings"
)

type Handler struct {
	store *database.TaskScore
}

func NewHandler(store *database.TaskScore) *Handler {
	return &Handler{store: store}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.GetAll()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting all tasks")
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

func (h *Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idString := pathParts[0]
	id, err := strconv.Atoi(idString)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Task ID")
	}

	task, err := h.store.GetByID(id)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	respondWithJSON(w, http.StatusOK, task)
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTaskInput
	//ссоздаем декодер для чтения из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if strings.TrimSpace(input.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "Missing title")
		return
	}

	if strings.TrimSpace(input.Description) == "" {
		respondWithError(w, http.StatusBadRequest, "Missing description")
		return
	}

	task, err := h.store.Create(input)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, task)
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idString := pathParts[0]
	id, err := strconv.Atoi(idString)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Task ID")
	}

	var input models.UpdateTaskInput
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Title != nil && len(*input.Title) == 0 {
		respondWithError(w, http.StatusBadRequest, "Missing title")
		return
	}

	task, err := h.store.Update(id, input)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, task)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idString := pathParts[0]
	id, err := strconv.Atoi(idString)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Task ID")
		return
	}

	err = h.store.Delete(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
