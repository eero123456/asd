package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"
	"todo/model"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func initDB() (err error) {

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	serverAddress := os.Getenv("DB_ADDRESS")
	database := os.Getenv("DB_DATABASE")

	connStr := fmt.Sprintf("%s:%s@%s/%s", user, password, serverAddress, database)

	DB, err = sql.Open("mysql", connStr)

	if err != nil {
		panic(err)
	}

	DB.SetConnMaxLifetime(time.Minute * 1)
	DB.SetMaxOpenConns(50)
	DB.SetMaxIdleConns(10)

	return
}

func test() (*[]model.Todo, error) {

	data := make([]model.Todo, 0)

	res, err := DB.Query("SELECT id,text,completed FROM todo")
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {

		var city model.Todo
		err := res.Scan(&city.ID, &city.Text, &city.Completed)

		if err != nil {
			return nil, err
		}

		data = append(data, city)
	}

	return &data, nil
}

func addUser() (userID uint64, err error) {
	row := DB.QueryRow("INSERT INTO user VALUES () RETURNING id")
	err = row.Scan(&userID)
	return
}

func deleteTodo(id string, userID uint64) error {

	stmt, err := DB.Prepare("DELETE FROM todo WHERE id=? AND user_id=?")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(id, userID)

	return err

}

func deleteTodos(ids []string, userID uint64) error {

	count := len(ids)
	vals := []interface{}{}

	bindStr := strings.Repeat("?,", count)
	bindStr = strings.TrimRight(bindStr, ",")

	for i := 0; i < count; i++ {

		vals = append(vals, ids[i])
	}

	sql := fmt.Sprintf("DELETE FROM todo WHERE user_id=%d AND id IN(%s)", userID, bindStr)
	//fmt.Println(sql)

	fmt.Println(sql)

	stmt, err := DB.Prepare(sql)

	if err != nil {
		return err
	}

	res, err := stmt.Exec(vals...)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	//res.RowsAffected()
	fmt.Println(res.RowsAffected())
	return nil

}

func deleteAllTodos(userID uint64) (err error) {

	stmt, err := DB.Prepare("DELETE FROM todo WHERE user_id=?;")

	if err != nil {
		return
	}

	_, err = stmt.Exec(userID)

	return

}

func saveTodo(text string, userID uint64) (inserted []model.Todo, err error) {

	stmt, err := DB.Prepare("INSERT INTO todo (text,user_id) VALUES (?,?) RETURNING id,text,completed;")
	if err != nil {
		return
	}
	var t model.Todo
	row := stmt.QueryRow(text, userID)
	err = row.Scan(&t.ID, &t.Text, &t.Completed)

	return []model.Todo{t}, err
}

func saveTodos(texts []string, userID uint64) ([]model.Todo, error) {

	count := len(texts)
	valueStrings := make([]string, 0, count)
	vals := []interface{}{}
	params := fmt.Sprintf("(?,%d)", userID)

	for i := 0; i < count; i++ {
		valueStrings = append(valueStrings, params)
		vals = append(vals, texts[i])
	}

	sql := fmt.Sprintf("INSERT INTO todo (text,user_id) VALUES %s RETURNING id,text,completed",
		strings.Join(valueStrings, ","))
	//fmt.Println(sql)
	stmt, err := DB.Prepare(sql)

	if err != nil {
		return nil, err
	}

	res, err := stmt.Query(vals...)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	defer res.Close()
	inserted := make([]model.Todo, 0)
	for res.Next() {

		var t model.Todo
		err := res.Scan(&t.ID, &t.Text, &t.Completed)

		if err != nil {
			return nil, err
		}

		inserted = append(inserted, t)
	}

	return inserted, nil

}

func getTodos(userID uint64) ([]model.Todo, error) {

	todos := make([]model.Todo, 0)

	res, err := DB.Query("SELECT id,text,completed FROM todo WHERE user_id=?", userID)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {

		var t model.Todo
		err := res.Scan(&t.ID, &t.Text, &t.Completed)

		if err != nil {
			return nil, err
		}

		todos = append(todos, t)
	}

	return todos, nil
}

func getEveryTodo() ([]model.Todo, error) {

	todos := make([]model.Todo, 0)

	res, err := DB.Query("SELECT id,text,completed FROM todo")
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {

		var t model.Todo
		err := res.Scan(&t.ID, &t.Text, &t.Completed)

		if err != nil {
			return nil, err
		}

		todos = append(todos, t)
	}

	return todos, nil
}

func getUserData(username string) (id int) {
	row := DB.QueryRow("SELECT id FROM user WHERE username=?", username)
	err := row.Scan(&id)
	if err != nil {
		return 0
	}
	return
}
