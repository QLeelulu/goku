package main

import (
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/mustache.goku"
    "log"
    "path"
    "runtime"
)

/**
 * Controller & Action
 */
var _ = goku.Controller("home").
    // home.index action
    Get("index", func(ctx *goku.HttpContext) goku.ActionResulter {
    return ctx.View(nil)
})

// routes
var routes []*goku.Route = []*goku.Route{
    // default controller and action route
    &goku.Route{
        Name:    "default",
        Pattern: "/{controller}/{action}/",
        Default: map[string]string{"controller": "home", "action": "index"},
    },
}

// server config
var config *goku.ServerConfig = &goku.ServerConfig{Addr: ":8080"}

func init() {
    // project root dir, this code can not put to main func
    _, filename, _, _ := runtime.Caller(1)
    config.RootDir = path.Dir(filename)
}

func main() {
    /**
     * set template engine to mustache
     */
    config.TemplateEnginer = mustache.NewMustacheTemplateEngine()

    config.LogLevel = goku.LOG_LEVEL_LOG

    rt := &goku.RouteTable{Routes: routes}
    s := goku.CreateServer(rt, nil, config)

    goku.Logger().Logln("Server start on", s.Addr)
    log.Fatal(s.ListenAndServe())
}
