package controllers

import (
    "github.com/qleelulu/goku"
    "github.com/qleelulu/goku/examples/todo/todo/models"
)

var _ = goku.Controller("home").
    Get("index", func(ctx *goku.HttpContext) goku.ActionResulter {
    todos, err := models.GetTodoLists()
    if err != nil {
        return ctx.Error(err)
    }
    return ctx.View(todos)
})
