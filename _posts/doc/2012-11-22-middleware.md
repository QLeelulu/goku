---
layout: post
title: "Middleware"
description: ""
category: doc
tags: []
---
{% include JB/setup %}

Sometimes we want to perform logic either before all the request or after all the request. To support this, Goku provides middlewares. Middleware provided a way let you to do something before or after a request.

![{{ ASSET_PATH }}goku-mf.png]({{ ASSET_PATH }}goku-mf.png)

## How To Create Middleware ##

To create a middleware, just implement the [Middlewarer](http://godoc.org/github.com/QLeelulu/goku#Middlewarer) interface.

{% highlight go %}
type Middlewarer interface {
    OnBeginRequest(ctx *HttpContext) (ActionResulter, error)
    OnBeginMvcHandle(ctx *HttpContext) (ActionResulter, error)
    OnEndMvcHandle(ctx *HttpContext) (ActionResulter, error)
    OnEndRequest(ctx *HttpContext) (ActionResulter, error)
}
{% endhighlight %}

middlewarer interface execute order: 

    1. OnBeginRequest
    2. OnBeginMvcHandle
    3.  -> {controller}
    4. OnEndMvcHandle
    5. OnEndRequest

notice that:

    OnBeginRequest & OnEndRequest: All requests will be through these
    OnBeginMvcHandle & OnEndMvcHandle: not matched route & static file are not through these

To add a middleware to the goku server, we can do this when we create goku server:

{% highlight go %}
func main() {
    rt := &goku.RouteTable{}
    middlewares := []goku.Middlewarer{
        new(middlewares.UserMiddleware),
        new(middlewares.LogMiddleware),
    }
    // add middlewares when we create server
    s := goku.CreateServer(rt, middlewares, serverConfig)
    log.Fatal(s.ListenAndServe())
}
{% endhighlight %}



