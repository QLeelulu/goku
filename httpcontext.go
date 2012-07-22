package goku

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
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

func (ctx *HttpContext) Get(name string) string {
    v, ok := ctx.RouteData.Get(name)
    if ok {
        return v
    }
    return ctx.Request.FormValue(name)
}

func (ctx *HttpContext) Header() http.Header {
    return ctx.responseWriter.Header()
}

func (ctx *HttpContext) SetHeader(key string, value string) {
    ctx.responseWriter.Header().Set(key, value)
    //ctx.responseHeaderCache.Set(key, value)
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

// render the view and return a ActionResulter
// it will find the view in these rules:
//      1. /{ViewPath}/{Controller}/{viewName}
//      2. /{ViewPath}/shared/{viewName}
func (ctx *HttpContext) Render(viewName string, viewModel interface{}) ActionResulter {
    vr := &ViewResult{
        ViewEngine: ctx.requestHandler.ViewEnginer,
        ViewData:   ctx.ViewData,
        ViewModel:  viewModel,
        ViewName:   viewName,
    }
    vr.Body = new(bytes.Buffer)
    vr.Headers = map[string]string{"Content-Type": "text/html"}
    return vr
}

// render the view and return a ActionResulter
// it will find the view in these rules:
//      1. /{ViewPath}/{Controller}/{action}
//      2. /{ViewPath}/shared/{action}
func (ctx *HttpContext) View(viewData interface{}) ActionResulter {
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

// return json string result
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
