package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/pat"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var rd *render.Render

type Todo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

var todoMap map[int]*Todo

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
}

func getTodoListHandler(w http.ResponseWriter, r *http.Request) {
	list := []*Todo{}
	for _, v := range todoMap {
		list = append(list, v)
	}
	rd.JSON(w, http.StatusOK, list)
}

func addTodoHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	id := len(todoMap)
	todoMap[id] = &Todo{id, name, false, time.Now()}
	rd.JSON(w, http.StatusOK, todoMap[id])
}

type Success struct {
	Success bool `json:"success"`
}

func removeTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	log.Print(id)
	if _, ok := todoMap[id]; ok {
		delete(todoMap, id)
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
}

func main() {
	todoMap = make(map[int]*Todo)
	rd = render.New()
	mux := pat.New()
	mux.HandleFunc("/", indexHandler)
	mux.Get("/todos", getTodoListHandler)
	mux.Post("/todos", addTodoHandler)
	mux.Delete("/todos/{id:[0-9]+}", removeTodoHandler)

	n := negroni.Classic()
	n.UseHandler(mux)
	http.ListenAndServe(":3000", n)
}
