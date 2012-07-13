package models

import (
    "github.com/qleelulu/goku"
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
    qi := goku.SqlQueryInfo{}
    qi.Order = "finished desc, id desc"
    list, err := db.GetStructs(Todo{}, qi)
    todos := make([]*Todo, len(list))
    for i, todo := range list {
        todos[i] = todo.(*Todo)
    }
    return todos, err
}
