package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"todo/model"
	"todo/templates"

	"github.com/joho/godotenv"
)

func main() {

	if err := loadEnv(); err != nil {
		fmt.Println("can't load env")
		os.Exit(1)
	}

	initDB()
	initSecureCookie()

	allTodos, _ = getEveryTodo()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	//http.HandleFunc("/", indexHandler)
	http.HandleFunc("/ping", withLogging(pingHandler))

	http.HandleFunc("/register", register)

	http.HandleFunc("/ssr/todo/deleteAll", ssrDeleteAllTodosHandler)
	http.HandleFunc("/ssr/todo/delete", ssrDeleteTodoHandler)
	http.HandleFunc("/ssr/todo", ssrTodoHandler)
	http.HandleFunc("/ssr/all", ssrTest)
	http.HandleFunc("/ssr/all-cached", ssrTest2)
	http.HandleFunc("/ssr/showBatch", ssrIndex)
	//http.HandleFunc("/ssr/cancelBatch", ssrIndex)

	http.HandleFunc("/ssr", ssrIndex)

	//http.HandleFunc("/ssr/todos2", requireAuth(ssrTodos2))

	http.HandleFunc("/todos", requireAuth(withLogging(todosHandler)))
	http.HandleFunc("/todo", requireAuth(withLogging(todoHandler)))

	http.ListenAndServe(":8080", nil)

}

var allTodos []model.Todo

func ssrIndex(w http.ResponseWriter, r *http.Request) {

	showBatchEditor := strings.Contains(r.URL.String(), "showBatch")

	userID, err := getUserID(r)
	if err != nil {
		handleNewUser(w)
		t := make([]model.Todo, 0)
		templates.Index(t, showBatchEditor).Render(r.Context(), w)
		return
	}

	todos, err := getTodos(userID)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	templates.Index(todos, showBatchEditor).Render(r.Context(), w)

}

func ssrTest(w http.ResponseWriter, r *http.Request) {

	showBatchEditor := strings.Contains(r.URL.String(), "showBatch")

	todos, err := getEveryTodo()

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	templates.Index(todos, showBatchEditor).Render(r.Context(), w)

}

func ssrTest2(w http.ResponseWriter, r *http.Request) {

	showBatchEditor := strings.Contains(r.URL.String(), "showBatch")

	templates.Index(allTodos, showBatchEditor).Render(r.Context(), w)

}

func ssrShowUserTodos(ctx context.Context, w http.ResponseWriter, userID uint64, showBatchEditor bool) {

	todos, err := getTodos(userID)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	templates.Index(todos, showBatchEditor).Render(ctx, w)
}

func loadEnv() error {
	return godotenv.Load(".env")

}

var todos []model.Todo

/*
func ssrTodos2(w http.ResponseWriter, r *http.Request) {

	//buf := &bytes.Buffer{}

	templates.UserTodos(todos).Render(r.Context(), w)
	//addLatencyInfo(w, r)
	//w.Write(buf.Bytes())

}

func ssrTodos(w http.ResponseWriter, r *http.Request) {

	todos, err := getEveryTodo()

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	templates.UserTodos(todos).Render(r.Context(), w)
	//buf := &bytes.Buffer{}

	//templates.UserTodos(todos).Render(r.Context(), buf)
	//addLatencyInfo(w, r)
	//w.Write(buf.Bytes())

}
*/
func pingHandler(w http.ResponseWriter, r *http.Request) {

	addLatencyInfo(w, r)
	w.Write([]byte("pong"))

}

func register(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	if r.Method == "OPTIONS" {

		w.WriteHeader(http.StatusOK)
		return
	}

	_, err := getUserID(r)
	if err != nil {
		handleNewUser(w)
		w.Write([]byte("token assigned"))
		return
	}

	w.Write([]byte("already ok"))

}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	_, err := getUserID(r)
	if err != nil {
		handleNewUser(w)
	}

}

func SetCookie(w http.ResponseWriter, token string) {

	expires := time.Now().AddDate(0, 0, 1)

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Secure:   false,
		Expires:  expires,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)

}
