package goku

import (
	"net/http"
	"bytes"
)

// http context
type HttpContext struct {
	Request        *http.Request       // http request
	ResponseWriter http.ResponseWriter // http response
	Method         string              // http method

	//self fields
	RouteData *RouteData     // route data
	Result    ActionResulter // action result
	Err       error          // process error
	User      string         // user name
	Canceled  bool           // cancel continue process the request and return

	// private fileds
	requestHandler *RequestHandler
}

func (ctx *HttpContext) Get(name string) string {
	v, ok := ctx.RouteData.Get(name)
	if ok {
		return v
	}
	return ctx.Request.FormValue(name)
}

func (ctx *HttpContext) SetHeader(key string, value string) {
	ctx.ResponseWriter.Header().Set(key, value)
}
func (ctx *HttpContext) GetHeader(key string) string {
	return ctx.Request.Header.Get(key)
}

func (ctx *HttpContext) ContentType(ctype string) {
	ctx.ResponseWriter.Header().Set("Content-Type", ctype)
}

func (ctx *HttpContext) Status(code int) {
	ctx.ResponseWriter.WriteHeader(code)
}

func (ctx *HttpContext) Write(b []byte) {
	ctx.ResponseWriter.Write(b)
}
func (ctx *HttpContext) WriteString(content string) {
	ctx.ResponseWriter.Write([]byte(content))
}

// func (ctx *HttpContext) Redirect(status int, url_ string) {
// 	ctx.ResponseWriter.Header().Set("Location", url_)
// 	ctx.ResponseWriter.WriteHeader(status)
// 	ctx.ResponseWriter.Write([]byte("Redirecting to: " + url_))
// }

// func (ctx *HttpContext) NotModified() {
// 	ctx.ResponseWriter.WriteHeader(304)
// }

// func (ctx *HttpContext) NotFound(message string) {
// 	ctx.ResponseWriter.WriteHeader(404)
// 	ctx.ResponseWriter.Write([]byte(message))
// }

// render the view and return a ActionResulter
func (ctx *HttpContext) Render(viewName string, viewData interface{}) ActionResulter {
	return &ViewResult{
		ViewEngine: ctx.requestHandler.ViewEnginer,
		ViewData:   viewData,
		ViewName:   viewName,
	}
}

// render the view and return a ActionResulter
func (ctx *HttpContext) View(viewName string, viewData interface{}) ActionResulter {
	return ctx.Render(viewName, viewData)
}

func (ctx *HttpContext) Redirect(url_ string) ActionResulter {
	return &ActionResult{
		StatusCode: 302,
		Headers:    map[string]string{"Content-Type": "text/html", "Location": url_},
		Body:       bytes.NewBufferString("Redirecting to: " + url_),
	}
}

func (ctx *HttpContext) RedirectPermanent(url_ string) ActionResulter {
	return &ActionResult{
		StatusCode: 301,
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
		StatusCode: 404,
		Headers:    map[string]string{"Content-Type": "text/html"},
		Body:       bytes.NewBufferString(message),
	}
}

// content not modified
func (ctx *HttpContext) NotModified(viewName string, viewData interface{}) ActionResulter {
	return &ActionResult{
		StatusCode: 304,
	}
}

func (ctx *HttpContext) Error(err interface{}) ActionResulter {
	var msg string
	switch t := err.(type) {
	case *string:
		msg = err.(string)
	case *error:
		msg = err.(error).Error()
	default:
		panic("wrong type: " + t.(string))
	}
	return &ActionResult{
		StatusCode: 500,
		Headers:    map[string]string{"Content-Type": "text/plain"},
		Body:       bytes.NewBufferString(msg),
	}

}

func (ctx *HttpContext) Raw(data string) ActionResulter {
	return &ActionResult{
		StatusCode: 500,
		Headers:    map[string]string{"Content-Type": "text/plain"},
		Body:       bytes.NewBufferString(data),
	}
}

// this.raw = function(data, contentType) {
//   return new ActionResult(data, {'Content-Type': contentType || 'text/plain'});
// };

// this.json = function(data) {
//   return new ActionResult(JSON.stringify(data), {'Content-Type': 'application/json'});
// };

// this.content = function(filename) {
//   return new ContentResult(filename);
// }
