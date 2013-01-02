---
layout: post
title: "Controller & Action"
description: ""
category: doc
tags: []
---
{% include JB/setup %}

Goku maps URLs to controllers that regiest to goku server. Controllers process incoming requests, handle user input and interactions, and execute appropriate application logic. A controller typically calls a separate view component to generate the HTML markup for the request.

## ControllerBuilder

We build controller with [ControllerBuilder](http://go.pkgdoc.org/github.com/QLeelulu/goku#ControllerBuilder)

{% highlight go %}
// a home controller builder
var cb *goku.ControllerBuilder = goku.Controller("home")
// add a "index" action to "home" controller for http get
cb.Get("index", func(ctx *goku.HttpContext) goku.ActionResulter {
    return ctx.Html("Hello World")
})
{% endhighlight %}

Use `goku.Controller("controller-name")` to get a controller builder. 
If the controller `"controller-name"` exist, it will just return the 
same controller builder.

## Action

User interaction with goku MVC applications is organized around controllers and action methods. The controller defines action methods. Controllers can include as many action methods as needed.

When a user enters a URL into the browser, the MVC application uses routing rules to parse the URL and to determine the path of the controller. The controller then determines the appropriate action method to handle the request. By default, the URL of a request is treated as a sub-path that includes the controller name followed by the action name. For example, if a user enters the URL `http://abc.com/Products/Categories`, the sub-path is `/Products/Categories`. The default routing rule treats "Products" as the controller. It treats "Categories" as the name of the action. Therefore, the routing rule invokes the Categories action of the Products controller in order to process the request. If the URL ends with `/Products/Detail/5`, the default routing rule treats `"Detail"` as the name of the action, and the Detail action of the Products controller is invoked to process the request. 

The following example shows a controller that has a `HelloWorld` action for http GET method.

{% highlight go %}
var _ = goku.Controller("my").
    Action("get", "HelloWorld", func(ctx *goku.HttpContext) goku.ActionResulter {
        return ctx.Html("Hello World")
    })
{% endhighlight %}

[ControllerBuilder](http://go.pkgdoc.org/github.com/QLeelulu/goku#ControllerBuilder) has some helper method: `Get`, `Post`, `Put`. So you can do it the same by this code:

{% highlight go %}
var _ = goku.Controller("my").
    Get("HelloWorld", func(ctx *goku.HttpContext) goku.ActionResulter {
        return ctx.Html("Hello World")
    })
{% endhighlight %}

If you set http method to `all`, the action will match all http method, but Priority is low.    
For example:

{% highlight go %}
var _ = goku.Controller("user").
    // action 1
    Post("edit", func(ctx *goku.HttpContext) goku.ActionResulter {
        return ctx.View(nil)
    }).
    // action 2
    Action("all", "edit", func(ctx *goku.HttpContext) goku.ActionResulter {
        return ctx.View(nil)
    })
{% endhighlight %}

A `Http POST /user/edit` request will process by `Post("edit", ...)`,     
but a `Http GET /user/edit` request will process by `Action("all", "edit", ...)`.


## ActionResult Return Type

Action methods return an instance of a struct that implement the 
[ActionResulter](http://go.pkgdoc.org/github.com/QLeelulu/goku#ActionResulter) interface. 
There are different action result types, depending on the task that the action method is performing. For example, the most common action is to call the `HttpContext.View` method. The View method returns an instance of the ViewResult struct, which is implement the ActionResulter interface.

The following table shows the built-in action result types and the action helper methods that return them.

Action Result | Helper Method | Description
--------------|---------------|------------
[ViewResult](http://godoc.org/github.com/QLeelulu/goku#ViewResult) | [ctx.View](http://godoc.org/github.com/QLeelulu/goku#HttpContext.View) | Renders a view as a Web page.
[ViewResult](http://godoc.org/github.com/QLeelulu/goku#ViewResult) | [ctx.RenderPartial](http://godoc.org/github.com/QLeelulu/goku#HttpContext.RenderPartial) | Renders a partial view, mean that not render the layout.
[ActionResulter](http://godoc.org/github.com/QLeelulu/goku#ActionResulter) | [ctx.Redirect](http://godoc.org/github.com/QLeelulu/goku#HttpContext.Redirect) | Redirects to another action method by using its URL.(`302`)
[ActionResulter](http://godoc.org/github.com/QLeelulu/goku#ActionResulter) | [ctx.RedirectPermanent](http://godoc.org/github.com/QLeelulu/goku#HttpContext.RedirectPermanent) | Redirects to another action method by using its URL.(`301`)
[ContentResult](http://godoc.org/github.com/QLeelulu/goku#ContentResult) | No | Returns binary output to write to the response.
[ActionResulter](http://godoc.org/github.com/QLeelulu/goku#ActionResulter) | [ctx.Json](http://godoc.org/github.com/QLeelulu/goku#HttpContext.Json) | Returns a serialized JSON object.
[ActionResulter](http://godoc.org/github.com/QLeelulu/goku#ActionResulter) | [ctx.Html](http://godoc.org/github.com/QLeelulu/goku#HttpContext.Html) | Returns string content, but set response Content-Type to text/html .
[ActionResulter](http://godoc.org/github.com/QLeelulu/goku#ActionResulter) | [ctx.Raw](http://godoc.org/github.com/QLeelulu/goku#HttpContext.Raw) | Returns string content, but set response Content-Type to text/plain .
[ActionResulter](http://godoc.org/github.com/QLeelulu/goku#ActionResulter) | [ctx.NotModified](http://godoc.org/github.com/QLeelulu/goku#HttpContext.NotModified) | Returns 304 not modified .
[ActionResulter](http://godoc.org/github.com/QLeelulu/goku#ActionResulter) | [ctx.NotFound](http://godoc.org/github.com/QLeelulu/goku#HttpContext.NotFound) | Returns 404 page not found .


## Custom ActionResult

To implement your own action result, you just need to implement the [ActionResulter](http://godoc.org/github.com/QLeelulu/goku#ActionResulter) interface, and return it in the action.

