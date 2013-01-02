---
layout: post
title: "HttpContext"
description: ""
category: doc
tags: []
---
{% include JB/setup %}

{% highlight go %}
type HttpContext struct {
    Request *http.Request // http request

    Method string // http method

    //self fields
    RouteData *RouteData             // route data
    ViewData  map[string]interface{} // view data for template
    Data      map[string]interface{} // data for httpcontex
    Result    ActionResulter         // action result
    Err       error                  // process error
    User      string                 // user name
    Canceled  bool                   // cancel continue process the request and return
    // contains filtered or unexported fields
}
{% endhighlight %}

If you want to shared data in a request, you can set it to `HttpContext.Data`, and then you can access it in the whole request.

For more detail, check [HttpContext doc](http://godoc.org/github.com/QLeelulu/goku#HttpContext).
