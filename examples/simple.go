package main

import (
    //"os"
    "github.com/qleelulu/goku"
    "log"
    "path"
    "runtime"
    "time"
)

// routes
var routes []*goku.Route = []*goku.Route{
    // static file route
    &goku.Route{
        Name:     "static",
        IsStatic: true,
        Pattern:  "/static/(.*)",
    },
    // default controller and action route
    &goku.Route{
        Name:       "default",
        Pattern:    "/{controller}/{action}/{id}",
        Default:    map[string]string{"controller": "home", "action": "index", "id": "0"},
        Constraint: map[string]string{"id": "\\d+"},
    },
}

// server config
var config *goku.ServerConfig = &goku.ServerConfig{
    Addr:           ":8888",
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
    //RootDir:        os.Getwd(),
    StaticPath: "static",
    ViewPath:   "views",
    Debug:      true,
}

func init() {
    /**
     * Controller & Action
     */
    goku.Controller("home").
        //Filters(new(TestControllerFilter)). // this filter is fot controller(all the actions)
        Get("index", func(ctx *goku.HttpContext) goku.ActionResulter {
        return ctx.Html("Hello World")
    })

    // project root dir
    _, filename, _, _ := runtime.Caller(1)
    config.RootDir = path.Dir(filename)
}

func main() {
    rt := &goku.RouteTable{Routes: routes}
    s := goku.CreateServer(rt, nil, config)
    goku.Logger().Logln("Server start on", s.Addr)
    log.Fatal(s.ListenAndServe())
}
