/**
 */

package goku

import (
	"net/http"
	"fmt"
	"time"
	"path"
	"bytes"
)

// all the config to the web server
type ServerConfig struct {
	Addr           string        // TCP address to listen on, ":http" if empty
	ReadTimeout    time.Duration // maximum duration before timing out read of the request
	WriteTimeout   time.Duration // maximum duration before timing out write of the response
	MaxHeaderBytes int           // maximum size of request headers, DefaultMaxHeaderBytes if 0

	RootDir    string // project root dir
	StaticPath string // static file dir, "static" if empty
	ViewPath   string // view file dir, "views" if empty

	Debug bool
}

// server inherit from http.Server
type Server struct {
	http.Server
}

// request handler, the main handler for all the requests
type RequestHandler struct {
	RouteTable        *RouteTable
	MiddlewareHandler MiddlewareHandler
	ServerConfig      *ServerConfig
	ViewEnginer       ViewEnginer
}

// implement the http.Handler interface
func (mh *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ctx *HttpContext
	ctx = mh.buildContext(w, r)
	var (
		ar  ActionResulter
		err error
	)
	ar, err = mh.execute(ctx)
	if err != nil {
		ar = ctx.Error(err)
	}
	if ar != nil {
		ar.ExecuteResult(ctx)
	}
	// response content was cached,
	// flush all the cached content to responsewriter
	ctx.flushToResponse()
}

// 你可以通过三种途径取消一个请求： 设置 ctx.Canceled = true , 返回一个ActionResulter或者一个错误
func (mh *RequestHandler) execute(ctx *HttpContext) (ar ActionResulter, err error) {
	// being request
	ar, err = mh.MiddlewareHandler.BeginRequest(ctx)
	if ctx.Canceled || err != nil || ar != nil {
		return
	}
	// match route
	routeData, ok := mh.RouteTable.Match(ctx.Request.URL.Path)
	if !ok {
		ar = ctx.NotFound("Page Not Found! No Route For The URL: " + ctx.Request.URL.Path)
		return
	}
	// static file route
	// return ContentResult
	if routeData.Route.IsStatic {
		sc := ctx.requestHandler.ServerConfig
		filePath := path.Join(sc.RootDir, sc.StaticPath, routeData.FilePath)
		fmt.Printf("fp: %s\n", filePath)
		ar = &ContentResult{
			FilePath: filePath,
		}
	} else {
		ctx.RouteData = routeData
		// parse form data before mvc handle
		ctx.Request.ParseForm()
		// begin mvc handle
		ar, err = mh.MiddlewareHandler.BeginMvcHandle(ctx)
		if ctx.Canceled || err != nil || ar != nil {
			return
		}
		// handle controller
		ar, err = mh.executeController(ctx)
		if ctx.Canceled || err != nil || ar != nil {
			return
		}
		// end mvc handle
		ar, err = mh.MiddlewareHandler.EndMvcHandle(ctx)
		if ctx.Canceled || err != nil || ar != nil {
			return
		}
	}
	// end request
	ar, err = mh.MiddlewareHandler.EndRequest(ctx)
	return
}

func (mh *RequestHandler) executeController(ctx *HttpContext) (ar ActionResulter, err error) {
	var ai *ActionInfo
	ai = defaultControllerFactory.GetAction(ctx.Method, ctx.RouteData.Controller, ctx.RouteData.Action)
	if ai == nil {
		ar = ctx.NotFound(fmt.Sprintf("No Action for Controlle:%s, Action:%s.",
			ctx.RouteData.Controller, ctx.RouteData.Action))
		return
	}
	// ing & ed filter's order is not the same
	ingFilters := append(ai.Controller.Filters, ai.Filters...)
	// action executing filter
	ar, err = runFilterActionExecuting(ctx, ingFilters)
	if ctx.Canceled || err != nil || ar != nil {
		return
	}
	// execute action
	var rar ActionResulter
	rar = ai.Handler(ctx)
	// action executed filter
	edFilters := append(ai.Filters, ai.Controller.Filters...)
	ar, err = runFilterActionExecuted(ctx, edFilters)
	if ctx.Canceled || err != nil || ar != nil {
		return
	}
	// resule executing filter
	ar, err = runFilterResultExecuting(ctx, ingFilters)
	if ctx.Canceled || err != nil || ar != nil {
		return
	}
	// execute action result
	rar.ExecuteResult(ctx)
	// result executed filter
	ar, err = runFilterResultExecuted(ctx, edFilters)
	return
}

func (mh *RequestHandler) buildContext(w http.ResponseWriter, r *http.Request) *HttpContext {
	//r.ParseForm()
	return &HttpContext{
		Request:              r,
		responseWriter:       w,
		Method:               r.Method,
		requestHandler:       mh,
		responseContentCache: new(bytes.Buffer),
		//responseHeaderCache: make(map[string]string),
	}
}

func (mh *RequestHandler) checkError(ctx *HttpContext, ar ActionResulter, err error) ActionResulter {
	if err != nil {
		return ctx.Error(err)
	}
	return ar
}

func CreateServer(routeTable *RouteTable, middlewares []Middlewarer, sc *ServerConfig) *Server {
	// if rootDir == "" {
	// 	panic("MvcServer: Root Dir must set")
	// }
	if routeTable == nil {
		panic("MvcServer: RouteTable is nil")
	}
	if routeTable.Routes == nil || len(routeTable.Routes) < 1 {
		panic("MvcServer: No Route in the RouteTable")
	}

	mh := &DefaultMiddlewareHandle{
		Middlewares: middlewares,
	}

	handler := &RequestHandler{
		RouteTable:        routeTable,
		MiddlewareHandler: mh,
		ServerConfig:      sc,
	}

	server := new(Server)
	server.Handler = handler
	server.Addr = sc.Addr
	server.ReadTimeout = sc.ReadTimeout
	server.WriteTimeout = sc.WriteTimeout
	server.MaxHeaderBytes = sc.MaxHeaderBytes
	return server
}
