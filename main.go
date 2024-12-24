package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Todo struct {
	ID    int
	Title string
	Done  bool
}

var (
	todos  = []Todo{}
	lastID = 0
	mu     sync.Mutex
	tpl    *template.Template
)

func main() {
	var err error
	tpl, err = template.ParseFiles("index.html")
	if err != nil {
		log.Fatal("Error loading template:", err)
	}

	http.HandleFunc("/", homePage)
	http.HandleFunc("/add", addTodos)
	http.HandleFunc("/delete", deleteTodo)
	fmt.Println("Server is running on port 8002")
	log.Fatal(http.ListenAndServe(":8002", nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	if err := tpl.ExecuteTemplate(w, "index.html", todos); err != nil {
		http.Error(w, "Unable to load page", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}

func addTodos(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		title := r.FormValue("title")
		mu.Lock()
		lastID++
		todos = append(todos, Todo{ID: lastID, Title: title, Done: false})
		mu.Unlock()
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid ID:", idStr)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, todo := range todos {
		if todo.ID == id {
			log.Println("Deleting:", todo) // Log the item to be deleted
			todos = append(todos[:i], todos[i+1:]...)
			break
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
