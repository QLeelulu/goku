package controllers

import (
    "github.com/QLeelulu/goku"
)

var _ = goku.Controller("home").
    Get("index", func(ctx *goku.HttpContext) goku.ActionResulter {
    return ctx.Redirect("/")
})
