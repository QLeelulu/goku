---
layout: post
category : doc
tags : [intro, beginner, tutorial]
---
{% include JB/setup %}

#About

goku is a Web Mvc Framework for golang, mostly like ASP.NET MVC.    
goku is simple and powerful.    
Base Features:  
+ mvc (Lightweight model) 
+ route 
+ multi template engine and layout 
+ simple database api 
+ form validation 
+ filter for controller or action 
+ middleware
+ and more

## [doc & api](http://go.pkgdoc.org/github.com/QLeelulu/goku)

> now goku is in the __`preview version, anything can be change`__.

##Installation

To install goku, simply run `go get github.com/QLeelulu/goku`.     
To use it in a program, use `import "github.com/QLeelulu/goku"`

## Example

Here is a simple example, or you can check [Intro](/doc/intro) for more detail.

{% highlight go %}
package main

import (
    "github.com/QLeelulu/goku"
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
     * project root dir
     */
    _, filename, _, _ := runtime.Caller(1)
    config.RootDir = path.Dir(filename)

    /**
     * Controller & Action
     */
    goku.Controller("home").
        Get("index", func(ctx *goku.HttpContext) goku.ActionResulter {
        return ctx.Html("Hello World")
    })

}

func main() {
    rt := &goku.RouteTable{Routes: routes}
    s := goku.CreateServer(rt, nil, config)
    goku.Logger().Logln("Server start on", s.Addr)
    log.Fatal(s.ListenAndServe())
}
{% endhighlight %}

## Examples

You can found some examples in the `github.com/QLeelulu/goku/examples` folder.    
To run example "todo" app, just:
    
{% highlight go %}
$ cd $GOROOT/src/pkg/github.com/QLeelulu/goku/examples/todo/
$ go run app.go
{% endhighlight %}

maybe you need run `todo.sql` first.

## Authors

 - [@QLeelulu](http://weibo.com/qleelulu)
 - waiting for you


## License

View the [LICENSE](https://github.com/senchalabs/connect/blob/master/LICENSE) file. 