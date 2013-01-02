
In the typical workflow of an MVC Web application, controller action methods handle an incoming Web request. These action methods use the incoming parameter values to execute application code and to retrieve or update data model objects from a database. The methods then select a view that renders a response to a browser.

## Rendering UI with Views ##

The following example shows how a view is rendered in a controller class.

{% highlight go %}
goku.Controller("home").
    Get("categories", func(ctx *goku.HttpContext) goku.ActionResulter {

    categories := northwind.GetCategories()
    ctx.ViewData["Title"] = "Goku Demo"
    return ctx.View(categories)
})
{% endhighlight %}

In the example, the parameter that is passed in the `View` method call is a slice of Category objects that are passed to the view. The [View](http://godoc.org/github.com/QLeelulu/goku#HttpContext.View) method calls the view engine, which uses the data in the slice to render to the view and to display it in the browser.

## View Engine ##
The MVC framework uses URL routing to determine which controller action to invoke, and the controller action then decides which views to render.

The `categories` action contains the following single line of code:

    return ctx.View(nil)

This line of code returns a view that must be located at the following path on your web server:

    /views/home/categories.html

The path to the view is inferred from the name of the controller and the name of the controller action.

If you prefer, you can be explicit about the view. The following line of code returns a view `all_categories` :

    ctx.Render("all_categories", nil);

When this line of code is executed, a view is returned from the following path:

    /views/home/all_categories.aspx

For example, goku's [ServerConfig.RootDir](http://godoc.org/github.com/QLeelulu/goku#ServerConfig) set to `"/data/www/"`, [ServerConfig.ViewPath](http://godoc.org/github.com/QLeelulu/goku#ServerConfig) set to `"views"`, and goku will find the view in these rule:

[ctx.View(ViewModel)](http://godoc.org/github.com/QLeelulu/goku#HttpContext.View): 

    1. /data/www/views/{controller}/{action}.html
    2. /data/www/views/shared/{action}.html
    
[ctx.Render("view_name", ViewModel)](http://godoc.org/github.com/QLeelulu/goku#HttpContext.Render): 
    
    1. /data/www/views/{controller}/view_name.html
    2. /data/www/views/shared/view_name.html

If the view name is start with `"/"`, like [ctx.Render("/home/index", ViewModel)](http://godoc.org/github.com/QLeelulu/goku#HttpContext.Render), it will find the view in this path:
    
    1. /data/www/views/home/index.html

and the layout will find in these rules:

    1. /{ViewPath}/{Controller}/{layout}
    2. /{ViewPath}/shared/{layout}

the default layout is `layout.html`, and you can change it in [ServerConfig.Layout](http://godoc.org/github.com/QLeelulu/goku#ServerConfig).

## Template Engine ##

The default template engine is go's `html/template`. To pass data to the template, you can do it in the action by two way:

{% highlight go %}
goku.Controller("home").
    Get("categories", func(ctx *goku.HttpContext) goku.ActionResulter {

    categories := northwind.GetCategories()
    ctx.ViewData["Title"] = "Goku Demo"
    return ctx.View(categories)
})
{% endhighlight %}

If you pass the data to the `ViewData` like `ctx.ViewData["Title"] = "Goku Demo"`, you can get it in the template by: 

{% capture text %}|.{ .Data.Title }.|{% endcapture %}{% include JB/liquid_raw %}

If you pass the data to the `ViewModel` like `return ctx.View(categories)`, you can get the categories in the template by:

{% capture text %}|.{ .Model }.|{% endcapture %}{% include JB/liquid_raw %}

So we can use it in the template like this:

{% capture text %}<!DOCTYPE HTML>
<html>
<head>
    <meta charset="UTF-8">
    <title>|.{ .Data.Title }.|</title>
</head>
<body>
    <h1>|.{ .Data.Title }.|</h1>
    <ul>
    |.{ range .Model }.|
    <li>|.{ .CategoryName }.|</li>
    |.{end}.|
    </ul>
</body>
</html>
{% endcapture %}
{% include JB/liquid_raw %}

## Layout ##

the layout will find in these rules:

    1. /{ViewPath}/{Controller}/{layout}
    2. /{ViewPath}/shared/{layout}

the default layout is `layout.html`, and you can change it in [ServerConfig.Layout](http://godoc.org/github.com/QLeelulu/goku#ServerConfig).

the layout.html:

{% capture text %}<!DOCTYPE HTML>
<html>
<head>
    <meta charset="UTF-8">
    <title>|.{ template "title" . }.|</title>
</head>
<body>
    <h1>Goku Demo</h1>
    |.{ template "body" . }.|
</body>
</html>
{% endcapture %}
{% include JB/liquid_raw %}

and then 