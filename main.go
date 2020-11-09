package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// Tutorial https://www.youtube.com/watch?v=pQAV8A9KLwk

type task struct {
	ID int `json:ID`
	Name string `json:Name`
	Content string `json:Content`
}

type allTasks []task

// Persistence
var tasks = allTasks{
	{
		ID:      1,
		Name:    "Task One",
		Content: "Some Content",
	},
}

func getTasks(w http.ResponseWriter, r *http.Request){
	w.Header().Set("content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request)  {
	var newTask task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprint(w, "Insert valid task")
	}
	json.Unmarshal(reqBody, &newTask)
	newTask.ID = len(tasks) +1
	tasks = append(tasks, newTask)

	w.Header().Set("content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func getTask(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	taskID, err := 	strconv.Atoi(vars["id"])

	if err != nil{
		fmt.Fprint(w, "Invalid ID")
	}
	for _, task := range tasks{
		if task.ID==taskID{
			w.Header().Set("content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	taskID, err := 	strconv.Atoi(vars["id"])

	if err != nil{
		fmt.Fprint(w, "Invalid ID")
	}
	for i, task := range tasks{
		if task.ID==taskID{
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Fprint(w, "Te task with ID ",taskID," Position ",i," was deleted")
		}
	}
}

func updateTask(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	taskID, err := 	strconv.Atoi(vars["id"])
	var newTask task
	if err != nil{
		fmt.Fprint(w, "Invalid ID")
	}
	for i, task := range tasks{
		if task.ID==taskID{
			tasks = append(tasks[:i], tasks[i+1:]...)
			newTask.ID=taskID
			newTask.Content = "Content updated"
			newTask.Name = "Updated task"
			tasks = append(tasks, newTask)
		}
	}
}

func indexRoute(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to my API")
}



func main()  {
	//Create new router
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")

	log.Fatal(http.ListenAndServe(":3000", router))
}


