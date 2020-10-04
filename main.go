package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Project struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Author   *Author `json:"author"`
}

type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var projects []Project

func getProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

func getProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range projects {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Project{})
}

func createProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var project Project
	_ = json.NewDecoder(r.Body).Decode(&project)
	project.ID = strconv.Itoa(rand.Intn(1000000)) //not safe
	projects = append(projects, project)
	json.NewEncoder(w).Encode(project)
}

func updateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range projects {
		if item.ID == params["id"] {
			projects = append(projects[:index], projects[index+1:]...)
			var project Project
			_ = json.NewDecoder(r.Body).Decode(&project)
			project.ID = params["id"]
			projects = append(projects, project)
			json.NewEncoder(w).Encode(project)
			return
		}
	}
	json.NewEncoder(w).Encode(projects)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range projects {
		if item.ID == params["id"] {
			projects = append(projects[:index], projects[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(projects)

}

func main() {

	router := mux.NewRouter()

	//mock data
	projects = append(projects,
		Project{
			ID:       "1",
			Name:     "go-api",
			Category: "programming",
			Author: &Author{
				Firstname: "Marcelo",
				Lastname:  "Colombo"}})

	projects = append(projects,
		Project{
			ID:       "2",
			Name:     "node-api",
			Category: "programming",
			Author: &Author{
				Firstname: "Marcelo",
				Lastname:  "Colombo"}})

	router.HandleFunc("/api/projects", getProjects).Methods("GET")
	router.HandleFunc("/api/projects/{id}", getProject).Methods("GET")
	router.HandleFunc("/api/projects", createProject).Methods("POST")
	router.HandleFunc("/api/projects/{id}", updateProject).Methods("PUT")
	router.HandleFunc("/api/projects/{id}", deleteProject).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
