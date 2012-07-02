#goku

  goku is a Web Mvc Framework for golang, mostly like ASP.NET MVC.

##Route
    
    var routes []*goku.Route = []*goku.Route{
        &goku.Route{
            Name:     "static",
            IsStatic: true,
            Pattern:  "/static/(.*)",
        },
        &goku.Route{
            Name:       "onstraint",
            Pattern:    "/{controller}/{action}/{id}",
            Default:    map[string]string{"controller": "home", "action": "index", "id": "0"},
            Constraint: map[string]string{"id": "\\d+"},
        },
        &goku.Route{
            Name:    "default",
            Pattern: "/{controller}/{action}/{id}",
            Default: map[string]string{"controller": "home", "action": "index", "id": "0"},
        },
    }

or

    rt := new(goku.RouteTable)
    rt.Static("staticFile", "/static/(.*)")
    rt.Map(
        "blog-page", // route name
        "/blog/p/{page}", // pattern
        map[string]string{"controller": "blog", "action": "page", "page": "0"}, // default
        map[string]string{"page": "\\d+"} // constraint
    )
    rt.Map(
        "default", // route name
        "/{controller}/{action}", // pattern
        map[string]string{"controller": "home", "action": "index"}, // default
    )

##Controller And Action

+ get  `/home/index` will call `exports.index_get`
+ post `/home/index` will call `exports.index`

        // controllers/home.js
        exports.index = function(fnNext){
            // this.req --> the httpRequest object
            // this.res --> the httpResponse object
            // this.routeData --> the route info
            // this.ar --> the actionresults method(view, raw, json, redirect, redirectPermanent, notFound, notModified, error)
            fnNext( this.ar.view({msg: 'hello world'}) );
        };
        exports.index_get = function(fnNext){
            fnNext( this.ar.view({msg: 'hello world'}) );
        };

  You must call the `fnNext` to continue handler the request.


##Action Filter

    var myFilter = function(){
        this.onControllerExecuting = function(ctx, fnNext){
            // ctx.req --> the httpRequest object
            // ctx.res --> the httpResponse object
            // ctx.routeData --> the route info
            // ctx.ar --> the actionresults method
            // you can log or check the auth here
            if(!checkAuth(ctx.req)){
                fnNext( this.ar.redirect('/login') );
            }else{
                fnNext();
            }
        };
        this.onControllerExecuted = function(ctx, fnNext){
            fnNext();
        };
        this.onActionExecuting = function(ctx, fnNext){
            fnNext();
        };
        this.onActionExecuted = function(ctx, fnNext){
            fnNext();
        };
        this.onResultExecuting = function(ctx, fnNext){
            fnNext();
        };
        this.onResultExecuted = function(ctx, fnNext){
            fnNext();
        };
    };
    
    // Add filter to the controller
    this.filters = [new myFilter()];
    
    exports.index = function(fnNext){
        fnNext( this.ar.view({msg: 'hello world'}) );
    };
    // Add filter to the action
    exports.index.filters = [new myFilter()];
  
  You must call the `fnNext` to continue handler the request. To end the request, you can put any actionResult to `fnNext`, just like `fnNext('end')`.
  
  Order of the filters execution is:
    1. `onControllerExecuting`
    2. `onActionExecuting`
    2.  the action execute
    3. `onResultExecuting`
    4.  the actionResult execute
    5. `onResultExecuted`
    6. `onActionExecuted`
    7. `onControllerExecuted`



##Middleware

        exports.beginRequest = function(ctx, fnNext){
            // ctx.req --> the httpRequest object
            // ctx.res --> the httpResponse object
            // ctx.ar --> the actionresults method
            // you can log or do other thing here
            fnNext();
        };
        
        exports.beginMvcHandler = function(ctx, fnNext){
            // ctx.req --> the httpRequest object
            // ctx.res --> the httpResponse object
            // ctx.routeData --> the route info
            // ctx.ar --> the actionresults method
            fnNext();
        };
        
        exports.endMvcHandler = function(ctx, fnNext){
            // ctx.req --> the httpRequest object
            // ctx.res --> the httpResponse object
            // ctx.routeData --> the route info
            // ctx.ar --> the actionresults method
            fnNext();
        };
        
        exports.endRequest = function(ctx, fnNext){
            // ctx.req --> the httpRequest object
            // ctx.res --> the httpResponse object
            // ctx.routeData --> the route info
            // ctx.ar --> the actionresults method
            fnNext();
        };
        
  Make sure call the `fnNext` to continue handler the request. To end the request, you can put any actionResult to `fnNext`, just like `fnNext('end')`.

  Order of the middleware event execution is:
    1. `beginRequest`
    2. `beginMvcHandler`(if not the static file request)
    2.    mvc handler (if not the static file request)
    3. `endMvcHandler`(if not the static file request)
    4. `endRequest`

##Request and Response object

    - `req.get`:  the querystring key-value object
    - `req.post`: the post form data kev-value object
    - `req.cookies`: the cookies key-value object
    
    - `res.cookies.set`: set the respose cookies: `res.cookies.set('name', 'value', {*options*})`. `options` is the same as cookie options.
    - `res.cookies.clear`: delete the cookie: `res.cookies.clear('name', {*options*})`

##ViewEngine

  The viewdata you pass to the view is in the `viewdata` object:
  `#{viewdata.***}` or `#{vd.***}`
  
  HtmlHelper: 
    + helper file put in `projectDir/helpers/`
    + eg. `projectDir/helpers/timeFormat.js`
    + then use in the view like this: `var tf = Helper('timeFormat');`  
  
  
## Authors

 - QLeelulu
 - fengmk2
 - waiting for you


## License

View the [LICENSE](https://github.com/senchalabs/connect/blob/master/LICENSE) file.