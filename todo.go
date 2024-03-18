package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todo/model"
)

func todoHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "OPTIONS":
		w.WriteHeader(http.StatusOK)

	case "POST":
		addTodo(w, r)

	case "DELETE":
		deleteTodoHandler(w, r)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(ctxUserID).(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("cant get userID"))
		return
	}

	var ids []string
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deleteTodos(ids, userID)

}

func addTodo(w http.ResponseWriter, r *http.Request) {

	setCORS(w)

	var newTodoTexts []string
	var err error
	if err = json.NewDecoder(r.Body).Decode(&newTodoTexts); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if len(newTodoTexts) > 1000 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(">1000 is too many"))
		return
	}

	userID, ok := r.Context().Value(ctxUserID).(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("cant get userID"))
		return
	}

	var inserted []model.Todo

	if len(newTodoTexts) == 1 {
		inserted, err = saveTodo(newTodoTexts[0], userID)
	} else {
		inserted, err = saveTodos(newTodoTexts, userID)
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	json.NewEncoder(w).Encode(inserted)

}

func todosHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "OPTIONS":
		w.WriteHeader(http.StatusOK)

	case "GET":
		getAllTodosHandler(w, r)

	case "DELETE":
		deleteAllTodosHandler(w, r)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func getAllTodosHandler(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(ctxUserID).(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("cant get userID"))
		return
	}

	data, err := getTodos(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	jsonSTR, err := json.Marshal(data)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	addLatencyInfo(w, r)

	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonSTR)
}

func deleteAllTodosHandler(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(ctxUserID).(uint64)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("cant get userID"))
		return
	}

	err := deleteAllTodos(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	addLatencyInfo(w, r)

}
