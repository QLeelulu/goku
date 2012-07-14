package controllers

import (
    "github.com/qleelulu/goku"
)

var _ = goku.Controller("home").
    Get("index", func(ctx *goku.HttpContext) goku.ActionResulter {
    return ctx.Redirect("/")
})
