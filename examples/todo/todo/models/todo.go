package models

import (
    "database/sql"
    "github.com/QLeelulu/goku"
    "time"
)

type Todo struct {
    Id       int
    Title    string
    Finished bool
    PostDate time.Time
}

func GetTodoLists() ([]*Todo, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()
    qi := goku.SqlQueryInfo{}
    qi.Order = "finished asc, id desc"
    list, err := db.GetStructs(Todo{}, qi)
    todos := make([]*Todo, len(list))
    for i, todo := range list {
        todos[i] = todo.(*Todo)
    }
    return todos, err
}

func GetTodo(id int) (Todo, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()

    var todo Todo = Todo{}
    err := db.GetStruct(&todo, "id=?", id)
    return todo, err
}

func SaveTodo(m map[string]interface{}) (sql.Result, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()
    r, err := db.Insert("todo", m)
    return r, err
}

func UpdateTodo(id int, m map[string]interface{}) (sql.Result, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()
    r, err := db.Update("todo", m, "id=?", id)
    return r, err
}

func DeleteTodo(id int) (sql.Result, error) {
    var db *goku.MysqlDB = GetDB()
    defer db.Close()
    r, err := db.Delete("todo", "id=?", id)
    return r, err
}
