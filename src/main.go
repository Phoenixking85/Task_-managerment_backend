package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Tasks struct {
	ID          string `json:"id"`
	TaskName    string `json:"task_name"`
	TaskDetails string `json:"task_details"`
	Date        string `json:"date"`
	Completed   bool   `json:"completed"`
}

var tasks []Tasks

func initSampleTasks() {
	task := Tasks{
		ID:          "1",
		TaskName:    "New project idea",
		TaskDetails: "Brainstorm new project ideas for Q3",
		Date:        "2023-11-25",
		Completed:   false,
	}
	tasks = append(tasks, task)

	task1 := Tasks{
		ID:          "2",
		TaskName:    "Weekly report",
		TaskDetails: "Prepare weekly status report for management",
		Date:        "2023-11-26",
		Completed:   true,
	}
	tasks = append(tasks, task1)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Task Manager API - Use endpoints to manage tasks")
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	taskID := params["id"]

	for _, task := range tasks {
		if task.ID == taskID {
			json.NewEncoder(w).Encode(task)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Task not found"})
}

func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task Tasks

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Invalid request format"})
		return
	}

	task.ID = strconv.Itoa(rand.Intn(10000))

	if task.Date == "" {
		task.Date = time.Now().Format("2006-01-02")
	}

	tasks = append(tasks, task)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	taskID := params["id"]

	for index, task := range tasks {
		if task.ID == taskID {

			tasks = append(tasks[:index], tasks[index+1:]...)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Task deleted"})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Task not found"})
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	taskID := params["id"]

	var updatedTask Tasks
	err := json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Invalid request format"})
		return
	}

	for index, task := range tasks {
		if task.ID == taskID {

			updatedTask.ID = taskID

			tasks[index] = updatedTask

			json.NewEncoder(w).Encode(updatedTask)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Task not found"})
}

func toggleTaskStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	taskID := params["id"]

	for index, task := range tasks {
		if task.ID == taskID {

			tasks[index].Completed = !tasks[index].Completed

			json.NewEncoder(w).Encode(tasks[index])
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Task not found"})
}

func setupServer() {

	router := mux.NewRouter()

	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	router.HandleFunc("/tasks/{id}/toggle", toggleTaskStatus).Methods("PATCH")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := corsHandler.Handler(router)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func main() {

	rand.Seed(time.Now().UnixNano())

	initSampleTasks()

	setupServer()
}
