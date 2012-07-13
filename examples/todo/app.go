package main

import (
    "github.com/qleelulu/goku"
    "github.com/qleelulu/goku/examples/todo/todo"
    _ "github.com/qleelulu/goku/examples/todo/todo/controllers" // notice this!! import controllers
    "log"
)

func main() {
    rt := &goku.RouteTable{Routes: todo.Routes}
    middlewares := []goku.Middlewarer{}
    s := goku.CreateServer(rt, middlewares, todo.Config)
    goku.Logger().Logln("Server start on", s.Addr)
    log.Fatal(s.ListenAndServe())
}
