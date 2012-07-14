package goku

import (
    "bytes"
    "fmt"
    "github.com/qleelulu/goku/utils"
    "html/template"
    "io"
    "path"
    "strings"
)

type ViewData struct {
    Data    map[string]interface{}
    Model   interface{}
    Globals map[string]interface{}
    Body    interface{} // if in layout template, this will set
}

// TemplateEnginer interface
type TemplateEnginer interface {
    // render the view with viewData and write to w
    Render(viewpath string, layoutPath string, viewData *ViewData, w io.Writer)
    // return whether the tempalte support layout
    SupportLayout() bool
    // template file ext name, default is ".html"
    Ext() string
}

// DefaultTemplateEngine
type DefaultTemplateEngine struct {
    ExtName       string
    UseCache      bool
    TemplateCache map[string]*template.Template
}

// template file ext name, default is ".html"
func (te *DefaultTemplateEngine) Ext() string {
    if te.ExtName == "" {
        return ".html"
    }
    return te.ExtName
}

// return whether the tempalte support layout
func (te *DefaultTemplateEngine) SupportLayout() bool {
    return true
}

func (te *DefaultTemplateEngine) Render(filepath string, layoutPath string, viewData *ViewData, wr io.Writer) {
    if layoutPath != "" {
        buf := new(bytes.Buffer)
        te.render(filepath, viewData, buf)
        viewData.Body = template.HTML(buf.String())
        te.render(layoutPath, viewData, wr)
    } else {
        te.render(filepath, viewData, wr)
    }
}

func (te *DefaultTemplateEngine) render(filepath string, viewData interface{}, wr io.Writer) {
    var tmpl *template.Template
    var err error
    if te.UseCache {
        tmpl = te.TemplateCache[filepath]
    }
    if tmpl == nil {
        tmpl, err = template.ParseFiles(filepath)
        if err != nil {
            panic("DefaultTemplateEngine.Render: parse template \"" + filepath + "\" error, " + err.Error())
        }
        if te.UseCache {
            te.TemplateCache[filepath] = tmpl
        }
    }

    err = tmpl.Execute(wr, viewData)
    if err != nil {
        panic(err)
    }
}

type ViewInfo struct {
    Controller, Action, View, Layout string
}

// ViewEnginer interface
// Render need a template engine
// so it may take a TemplateEnginer
type ViewEnginer interface {
    Render(vi *ViewInfo, viewData *ViewData, wr io.Writer)
    // find the view and layout
    // if template engine not suppot layout, just return empty string
    Lookup(vi *ViewInfo) (viewPath string, layoutPath string)
}

// DefaultViewEngine
type DefaultViewEngine struct {
    RootDir               string // view's root dir, must set
    Layout                string // template layout name, default is "layout"
    ViewLocationFormats   []string
    LayoutLocationFormats []string
    TemplateEngine        TemplateEnginer
    UseCache              bool              // whether cache the viewfile
    Caches                map[string]string // controller & action & view to the real-file-path cache
}

func (ve *DefaultViewEngine) Lookup(vi *ViewInfo) (viewPath string, layoutPath string) {
    viewPath = ve.lookup(vi, false)
    if ve.TemplateEngine.SupportLayout() {
        layoutPath = ve.lookup(vi, true)
    }
    return
}

func (ve *DefaultViewEngine) lookup(vi *ViewInfo, isLayout bool) string {
    var viewName, cacheKey string
    var locas []string
    if isLayout {
        viewName = vi.Layout
        if vi.Layout == "" {
            viewName = ve.Layout // default layout
        } else {
            viewName = vi.Layout
        }
        if viewName == "" {
            return ""
        }
        cacheKey = vi.Controller + "_layout_" + viewName
        locas = ve.LayoutLocationFormats
    } else {
        viewName = vi.View
        if viewName == "" {
            viewName = vi.Action
        }
        cacheKey = vi.Controller + "_" + viewName
        locas = ve.ViewLocationFormats
    }
    viewName = viewName + ve.TemplateEngine.Ext()
    if ve.UseCache {
        if v, ok := ve.Caches[cacheKey]; ok {
            return v
        }
    }
    lookPaths := make([]string, 0, 3)
    for _, format := range locas {
        viewPath := strings.Replace(format, "{1}", vi.Controller, 1)
        viewPath = strings.Replace(viewPath, "{0}", viewName, 1)
        viewPath = path.Join(ve.RootDir, viewPath)
        if ok, _ := utils.FileExists(viewPath); ok {
            ve.Caches[cacheKey] = viewPath
            return viewPath
        }
        lookPaths = append(lookPaths, viewPath)
    }
    if !isLayout {
        panic(fmt.Sprintf("DefaultViewEngine: can't find the view for {controller: %s, action: %s, view: %s}, look up paths: %s",
            vi.Controller, vi.Action, vi.View, lookPaths))
    }
    return ""
}

func (ve *DefaultViewEngine) Render(vi *ViewInfo, viewData *ViewData, wr io.Writer) {
    viewFile, layoutFile := ve.Lookup(vi)
    ve.TemplateEngine.Render(viewFile, layoutFile, viewData, wr)
}

// create a default ViewEnginer
// if TemplateEnginer is nil, will use the default template engine
// some default value:
// 		+ Layout: "layout"
// 		+ ViewLocationFormats:   []string{"{1}/{0}", "shared/{0}"} , {1} is controller, {0} is action or a viewName
// 		+ LayoutLocationFormats: []string{"{1}/{0}", "shared/{0}"}
func CreateDefaultViewEngine(viewDir string, te TemplateEnginer, layout string, useCache bool) *DefaultViewEngine {
    if viewDir == "" {
        panic("CreateDefaultViewEngine: viewDir can not be empty.")
    }
    dve := &DefaultViewEngine{
        RootDir:        viewDir,
        Layout:         layout,
        TemplateEngine: te,
        UseCache:       useCache,
        Caches:         make(map[string]string),
    }
    // default location paths
    dve.ViewLocationFormats = []string{
        "{1}/{0}",
        "shared/{0}",
    }
    dve.LayoutLocationFormats = []string{
        "{1}/{0}",
        "shared/{0}",
    }
    if dve.Layout == "" {
        dve.Layout = "layout"
    }
    // DefaultTemplateEngine
    if dve.TemplateEngine == nil {
        dve.TemplateEngine = &DefaultTemplateEngine{
            UseCache:      useCache,
            TemplateCache: make(map[string]*template.Template),
        }
    }
    return dve
}

var globalViewData map[string]interface{} = make(map[string]interface{})

// add a view data to the global,
// that all the view can use it
// by {{.Global.key}}
func SddGlobalViewData(key string, val interface{}) {
    globalViewData[key] = val
}
