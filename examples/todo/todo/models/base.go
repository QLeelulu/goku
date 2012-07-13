package models

import (
    // _ "code.google.com/p/go-mysql-driver/mysql"
    "github.com/qleelulu/goku"
    "github.com/qleelulu/goku/examples/todo/todo"
    _ "github.com/ziutek/mymysql/godrv"
)

func GetDB() *goku.MysqlDB {
    db, err := goku.OpenMysql(todo.DATABASE_Driver, todo.DATABASE_DSN)
    if err != nil {
        panic(err.Error())
    }
    return db
}
