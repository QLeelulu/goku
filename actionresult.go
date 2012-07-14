package goku

import (
    "bytes"
    "io"
    "net/http"
    //"fmt"
)

type ActionResulter interface {
    ExecuteResult(ctx *HttpContext)
}

type ActionResult struct {
    StatusCode int
    Headers    map[string]string
    Body       *bytes.Buffer
}

func (ar *ActionResult) ExecuteResult(ctx *HttpContext) {
    if ar.Headers != nil {
        for k, v := range ar.Headers {
            ctx.SetHeader(k, v)
        }
    }
    // if ar.StatusCode == 0 {
    // 	ar.StatusCode = 200
    // }
    ctx.Status(ar.StatusCode)
    if ar.Body != nil && ar.Body.Len() > 0 {
        // TODO: which way is the fastest ?
        //ctx.Write(ar.Body.Bytes())
        //ar.Body.WriteTo(ctx.responseWriter)
        ctx.WriteBuffer(ar.Body)
    }
}

type ViewResult struct {
    ActionResult

    ViewEngine ViewEnginer
    ViewData   map[string]interface{}
    ViewModel  interface{}
    ViewName   string
}

func (vr *ViewResult) Render(ctx *HttpContext, wr io.Writer) {
    if vr.ViewEngine == nil {
        vr.ViewEngine = ctx.requestHandler.ViewEnginer
    }
    vi := &ViewInfo{
        Controller: ctx.RouteData.Controller,
        Action:     ctx.RouteData.Action,
        View:       vr.ViewName,
        Layout:     "",
    }
    viewData := &ViewData{
        Data:    vr.ViewData,
        Model:   vr.ViewModel,
        Globals: globalViewData,
    }
    vr.ViewEngine.Render(vi, viewData, wr)
}

func (vr *ViewResult) ExecuteResult(ctx *HttpContext) {
    vr.Render(ctx, vr.Body)
    vr.ActionResult.ExecuteResult(ctx)
}

type ContentResult struct {
    FilePath string
}

func (cr *ContentResult) ExecuteResult(ctx *HttpContext) {
    http.ServeFile(ctx, ctx.Request, cr.FilePath)
}
