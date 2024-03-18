package main

import (
	"fmt"
	"net/http"
	"strings"
)

func ssrTodoHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "POST":
		ssrAddTodo(w, r)

	/*
		case "DELETE":
			deleteTodo(w, r)
	*/
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func ssrDeleteAllTodosHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	err = deleteAllTodos(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, r, "/ssr", http.StatusSeeOther)

}

func ssrDeleteTodoHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	values := r.URL.Query()

	id := values.Get("id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	err = deleteTodo(id, userID)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
	}

	http.Redirect(w, r, "/ssr", http.StatusSeeOther)
}

func ssrAddTodo(w http.ResponseWriter, r *http.Request) {

	userID, err := getUserID(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	r.ParseForm()
	fmt.Println(r.Form)

	text := r.Form.Get("newTodo")
	if len(text) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("empty todo"))
		return
	}

	if !r.Form.Has("isBatch") {
		_, err = saveTodo(text, userID)
	} else {

		newTodoTexts := strings.Split(text, "\n")
		if len(newTodoTexts) > 1000 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(">1000 is too many"))
			return
		}

		_, err = saveTodos(newTodoTexts, userID)

	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, r, "/ssr", http.StatusSeeOther)
}
