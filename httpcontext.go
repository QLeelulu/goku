package goku

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "path"
)

// http context
type HttpContext struct {
    Request        *http.Request       // http request
    responseWriter http.ResponseWriter // http response
    Method         string              // http method

    //self fields
    RouteData *RouteData             // route data
    ViewData  map[string]interface{} // view data for template
    Data      map[string]interface{} // data for httpcontex
    Result    ActionResulter         // action result
    Err       error                  // process error
    User      string                 // user name
    Canceled  bool                   // cancel continue process the request and return

    // private fileds
    requestHandler       *RequestHandler
    responseContentCache *bytes.Buffer // cache response content, will write at end request
    responseStatusCode   int           // cache response status code, will write at end request
    //responseHeaderCache  Header        // cache response header, will write at end request
}

func (ctx *HttpContext) flushToResponse() {
    // if len(ctx.responseHeaderCache) > 0 {
    // 	for k, v := range ctx.responseHeaderCache {
    // 		ctx.responseWriter.Header().Set(key, value)
    // 	}
    // }
    if ctx.responseStatusCode > 0 {
        ctx.responseWriter.WriteHeader(ctx.responseStatusCode)
    }
    if ctx.responseContentCache.Len() > 0 {
        ctx.responseContentCache.WriteTo(ctx.responseWriter)
    }
}

// Try not to use this unless you know exactly what you are doing
func (ctx *HttpContext) ResponseWriter() http.ResponseWriter {
    return ctx.responseWriter
}

func (ctx *HttpContext) RootDir() string {
    return ctx.requestHandler.ServerConfig.RootDir
}

func (ctx *HttpContext) StaticPath() string {
    return path.Join(ctx.RootDir(), ctx.requestHandler.ServerConfig.StaticPath)
}

func (ctx *HttpContext) ViewPath() string {
    return path.Join(ctx.RootDir(), ctx.requestHandler.ServerConfig.ViewPath)
}

// get the requert param, 
// get from RouteData first, 
// if no, get from Requet.FormValue
func (ctx *HttpContext) Get(name string) string {
    v, ok := ctx.RouteData.Get(name)
    if ok {
        return v
    }
    return ctx.Request.FormValue(name)
}

// Header gets the response header
func (ctx *HttpContext) Header() http.Header {
    return ctx.responseWriter.Header()
}

// set the response header
func (ctx *HttpContext) SetHeader(key string, value string) {
    ctx.responseWriter.Header().Set(key, value)
    //ctx.responseHeaderCache.Set(key, value)
}

// AddHeader adds response header
func (ctx *HttpContext) AddHeader(key string, value string) {
    ctx.responseWriter.Header().Add(key, value)
}

// set response cookie header
func (ctx *HttpContext) SetCookie(cookie *http.Cookie) {
    ctx.responseWriter.Header().Add("Set-Cookie", cookie.String())
}

func (ctx *HttpContext) GetHeader(key string) string {
    return ctx.Request.Header.Get(key)
}

func (ctx *HttpContext) ContentType(ctype string) {
    ctx.responseWriter.Header().Set("Content-Type", ctype)
    //ctx.responseHeaderCache["Content-Type"] = ctype
}

func (ctx *HttpContext) Status(code int) {
    //ctx.responseWriter.WriteHeader(code)
    ctx.responseStatusCode = code
}

func (ctx *HttpContext) Write(b []byte) (int, error) {
    //return ctx.ResponseWriter.Write(b)
    return ctx.responseContentCache.Write(b)
}

func (ctx *HttpContext) WriteBuffer(bf *bytes.Buffer) {
    //bf.WriteTo(ctx.ResponseWriter)
    bf.WriteTo(ctx.responseContentCache)
}

func (ctx *HttpContext) WriteString(content string) {
    //ctx.ResponseWriter.Write([]byte(content))
    ctx.responseContentCache.Write([]byte(content))
}

func (ctx *HttpContext) WriteHeader(code int) {
    //ctx.responseWriter.WriteHeader(code)
    ctx.responseStatusCode = code
}

// IsAjax gets whether the request is by ajax
func (ctx *HttpContext) IsAjax() bool {
    return ctx.GetHeader("X-Requested-With") == "XMLHttpRequest"
}

// render the view and return a *ViewResult.
// it will find the view in these rules:
//      1. /{ViewPath}/{Controller}/{viewName}
//      2. /{ViewPath}/shared/{viewName}
// if viewName start with '/',
// it will find the view direct by viewpath:
//      1. /{ViewPath}/{viewName}
func (ctx *HttpContext) Render(viewName string, viewModel interface{}) *ViewResult {
    return ctx.rederView(viewName, "", viewModel, false)
}

// RenderWithLayout renders the view and return a *ViewResult
// it will find the view in these rules:
//      1. /{ViewPath}/{Controller}/{viewName}
//      2. /{ViewPath}/shared/{viewName}
func (ctx *HttpContext) RenderWithLayout(viewName, layout string, viewModel interface{}) *ViewResult {
    return ctx.rederView(viewName, layout, viewModel, false)
}

// RenderPartial renders a Partial view and return a *ViewResult.
// this is not use layout.
// it will find the view in these rules:
//      1. /{ViewPath}/{Controller}/{viewName}
//      2. /{ViewPath}/shared/{viewName}
func (ctx *HttpContext) RenderPartial(viewName string, viewModel interface{}) *ViewResult {
    return ctx.rederView(viewName, "", viewModel, true)
}

func (ctx *HttpContext) rederView(viewName, layout string, viewModel interface{}, isPartial bool) *ViewResult {
    vr := &ViewResult{
        ViewEngine:     ctx.requestHandler.ViewEnginer,
        TemplateEngine: ctx.requestHandler.TemplateEnginer,
        ViewData:       ctx.ViewData,
        ViewModel:      viewModel,
        ViewName:       viewName,
        Layout:         layout,
        IsPartial:      isPartial,
    }
    vr.Body = new(bytes.Buffer)
    vr.Headers = map[string]string{"Content-Type": "text/html"}
    return vr
}

// View renders the view and return a *ViewResult
// it will find the view in these rules:
//      1. /{ViewPath}/{Controller}/{action}
//      2. /{ViewPath}/shared/{action}
func (ctx *HttpContext) View(viewData interface{}) *ViewResult {
    return ctx.Render("", viewData)
}

func (ctx *HttpContext) Redirect(url_ string) ActionResulter {
    return &ActionResult{
        StatusCode: http.StatusFound,
        Headers:    map[string]string{"Content-Type": "text/html", "Location": url_},
        Body:       bytes.NewBufferString("Redirecting to: " + url_),
    }
}

func (ctx *HttpContext) RedirectPermanent(url_ string) ActionResulter {
    return &ActionResult{
        StatusCode: http.StatusMovedPermanently,
        Headers:    map[string]string{"Content-Type": "text/html", "Location": url_},
        Body:       bytes.NewBufferString("Redirecting to: " + url_),
    }
}

// page not found
func (ctx *HttpContext) NotFound(message string) ActionResulter {
    if message == "" {
        message = "Page Not Found!"
    }
    return &ActionResult{
        StatusCode: http.StatusNotFound,
        Headers:    map[string]string{"Content-Type": "text/html"},
        Body:       bytes.NewBufferString(message),
    }
}

// content not modified
func (ctx *HttpContext) NotModified() ActionResulter {
    return &ActionResult{
        StatusCode: http.StatusNotModified,
    }
}

func (ctx *HttpContext) Error(err interface{}) ActionResulter {
    msg := fmt.Sprintf("%v", err)
    return &ActionResult{
        StatusCode: http.StatusInternalServerError,
        Headers:    map[string]string{"Content-Type": "text/plain"},
        Body:       bytes.NewBufferString(msg),
    }

}

func (ctx *HttpContext) Raw(data string) ActionResulter {
    return &ActionResult{
        StatusCode: http.StatusOK,
        Headers:    map[string]string{"Content-Type": "text/plain"},
        Body:       bytes.NewBufferString(data),
    }
}

func (ctx *HttpContext) Html(data string) ActionResulter {
    return &ActionResult{
        StatusCode: http.StatusOK,
        Headers:    map[string]string{"Content-Type": "text/html"},
        Body:       bytes.NewBufferString(data),
    }
}

// Json returns json string result
// ctx.Json(obj) or ctx.Json(obj, "text/html")
func (ctx *HttpContext) Json(data interface{}, contentType ...string) ActionResulter {
    var ct string
    if len(contentType) == 1 {
        ct = contentType[0]
    } else {
        ct = "text/javascript"
    }
    ar := &ActionResult{
        StatusCode: http.StatusOK,
        Headers:    map[string]string{"Content-Type": ct},
        Body:       new(bytes.Buffer),
    }
    ec := json.NewEncoder(ar.Body)
    ec.Encode(data)
    return ar
}

// this.content = function(filename) {
//   return new ContentResult(filename);
// }
