package goku

import (
    //"fmt"
    "html/template"
    "net/http"
    "os"
    "path"
    "runtime"
)

type devErrorContext struct {
    ShowDetail bool
    Request    *http.Request
    Err        string
    StatusCode int
    Stack      string

    OsEnviron      []string
    GoRoot         string
    GoNumGoroutine int
    GoVersion      string
}

type devErrorHanller struct {
    view            string
    TemplateEnginer TemplateEnginer
}

func (eh *devErrorHanller) showErrorInfo(ctx *HttpContext, err string, statusCode int, showDetail bool, stack string) {
    ec := &devErrorContext{
        ShowDetail: showDetail,
        Request:    ctx.Request,
        Err:        err,
        StatusCode: statusCode,
        GoVersion:  runtime.Version(),
    }
    if showDetail {
        ec.Stack = stack
        ec.OsEnviron = os.Environ()
        ec.GoRoot = runtime.GOROOT()
        ec.GoNumGoroutine = runtime.NumGoroutine()
    }

    eh.TemplateEnginer.Render(eh.view, "", ec, ctx.responseContentCache)
}

func createDevErrorHandler() *devErrorHanller {
    //pwd, _ := os.Getwd()
    _, filename, _, _ := runtime.Caller(1)
    pwd := path.Dir(filename)
    eh := &devErrorHanller{
        view: path.Join(pwd, "views/error.html"),
        TemplateEnginer: &DefaultTemplateEngine{
            UseCache:      false, // true
            TemplateCache: make(map[string]*template.Template),
        },
    }
    return eh
}

var devErrorHanlle *devErrorHanller = createDevErrorHandler()

type devErrorResult struct {
    StatusCode int
    Err        string
    ShowDetail bool
    Stack      string
}

func (er *devErrorResult) ExecuteResult(ctx *HttpContext) {
    ctx.responseContentCache.Reset()
    ctx.Status(er.StatusCode)
    devErrorHanlle.showErrorInfo(ctx, er.Err, er.StatusCode, er.ShowDetail, er.Stack)
}
