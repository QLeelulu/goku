package models

import (
    _ "code.google.com/p/go-mysql-driver/mysql"
    "../../../../../goku"
    "../../todo"
)

func GetDB() *goku.MysqlDB {
    db, err := OpenMysql(todo.DATABASE_Driver, todo.DATABASE_DSN)
    if err != nil {
        panic(err.Error())
    }
    return db
}
