package controllers

import (
    "github.com/qleelulu/goku"
    // "github.com/qleelulu/goku/examples/todo/todo/models"
)

var _ = goku.Controller("user").
    Get("profile", func(ctx *goku.HttpContext) goku.ActionResulter {
    return ctx.Html("")
})
