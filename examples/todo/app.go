package main

import (
    "../../../goku"
    "./todo"
    _ "./todo/controllers"
    _ "runtime"
    _ "path"
    "log"
)

func main() {
    rt := &goku.RouteTable{Routes: todo.Routes}
    middlewares := []goku.Middlewarer{}
    s := goku.CreateServer(rt, middlewares, todo.Config)
    log.Fatal(s.ListenAndServe())
}
