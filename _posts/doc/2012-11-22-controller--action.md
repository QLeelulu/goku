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
ViewResult | View | Renders a view as a Web page.
PartialViewResult | PartialView | Renders a partial view, which defines a section of a view that can be rendered inside another view.
RedirectResult | Redirect | Redirects to another action method by using its URL.
RedirectToRouteResult | RedirectToAction <br/> RedirectToRoute | Redirects to another action method.
ContentResult | Content | Returns a user-defined content type.
JsonResult | Json | Returns a serialized JSON object.
JavaScriptResult | JavaScript | Returns a script that can be executed on the client.
FileResult | File | Returns binary output to write to the response.
EmptyResult | (None) | Represents a return value that is used if the action method must return a null result (void).


## Custom ActionResult


