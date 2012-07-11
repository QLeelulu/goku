package controllers

import (
    "../../../../../goku"
)

var _ = goku.Controller("home").
    Get("index", func(ctx *goku.HttpContext) goku.ActionResulter {
    return ctx.View(nil)
})
