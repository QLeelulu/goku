package goku

import (
    "fmt"
    "github.com/QLeelulu/goku/utils"
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
    if te.SupportLayout() && layoutPath != "" {
        te.render([]string{layoutPath, filepath}, viewData, wr)
    } else {
        te.render([]string{filepath}, viewData, wr)
    }
}

func (te *DefaultTemplateEngine) render(filepaths []string, viewData interface{}, wr io.Writer) {
    var tmpl *template.Template
    var err error
    cacheKey := strings.Join(filepaths, "_")
    if te.UseCache {
        tmpl = te.TemplateCache[cacheKey]
    }
    if tmpl == nil {
        tmpl, err = template.ParseFiles(filepaths...)
        if err != nil {
            panic("DefaultTemplateEngine.Render: parse template \"" + strings.Join(filepaths, ", ") + "\" error, " + err.Error())
        }
        if te.UseCache {
            te.TemplateCache[cacheKey] = tmpl
        }
    }

    err = tmpl.Execute(wr, viewData)
    if err != nil {
        panic(err)
    }
}

type ViewInfo struct {
    Controller, Action, View, Layout string
    IsPartial                        bool
}

// ViewEnginer interface.
// For how to find the view file.
type ViewEnginer interface {
    // find the view and layout
    // if template engine not suppot layout, just return empty string
    FindView(vi *ViewInfo) (viewPath string, layoutPath string)
}

// DefaultViewEngine
type DefaultViewEngine struct {
    ExtName               string // template file ext name, default is ".html"
    RootDir               string // view's root dir, must set
    Layout                string // template layout name, default is "layout"
    ViewLocationFormats   []string
    LayoutLocationFormats []string
    UseCache              bool              // whether cache the viewfile
    Caches                map[string]string // controller & action & view to the real-file-path cache
}

func (ve *DefaultViewEngine) FindView(vi *ViewInfo) (viewPath string, layoutPath string) {
    viewPath = ve.lookup(vi, false)
    if !vi.IsPartial {
        layoutPath = ve.lookup(vi, true)
    }
    return
}

func (ve *DefaultViewEngine) lookup(vi *ViewInfo, isLayout bool) string {
    var viewName, cacheKey string
    var locas []string
    if !vi.IsPartial && isLayout {
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
    viewName = viewName + ve.ExtName
    if ve.UseCache {
        if v, ok := ve.Caches[cacheKey]; ok {
            return v
        }
    }
    lookPaths := make([]string, 0, 3)
    // Absolute path, 
    // direct use viewpath
    if viewName[0] == '/' {
        viewPath := path.Join(ve.RootDir, viewName)
        if ok, _ := utils.FileExists(viewPath); ok {
            ve.Caches[cacheKey] = viewPath
            return viewPath
        }
        lookPaths = append(lookPaths, viewPath)
    } else {
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
    }
    if !isLayout {
        panic(fmt.Sprintf("DefaultViewEngine: can't find the view for {controller: %s, action: %s, view: %s}, look up paths: %s",
            vi.Controller, vi.Action, vi.View, lookPaths))
    }
    return ""
}

// create a default ViewEnginer. 
// some default value:
// 		+ Layout: "layout"
// 		+ ViewLocationFormats:   []string{"{1}/{0}", "shared/{0}"} , {1} is controller, {0} is action or a viewName
// 		+ LayoutLocationFormats: []string{"{1}/{0}", "shared/{0}"}
func CreateDefaultViewEngine(viewDir, layout, extName string, useCache bool) *DefaultViewEngine {
    if viewDir == "" {
        panic("CreateDefaultViewEngine: viewDir can not be empty.")
    }
    dve := &DefaultViewEngine{
        ExtName:  extName,
        RootDir:  viewDir,
        Layout:   layout,
        UseCache: useCache,
        Caches:   make(map[string]string),
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
    return dve
}

// create a default TemplateEnginer.
func CreateDefaultTemplateEngine(useCache bool) *DefaultTemplateEngine {

    te := &DefaultTemplateEngine{
        UseCache:      useCache,
        TemplateCache: make(map[string]*template.Template),
    }
    return te
}

var globalViewData map[string]interface{} = make(map[string]interface{})

// add a view data to the global,
// that all the view can use it
// by {{.Global.key}}
func SetGlobalViewData(key string, val interface{}) {
    globalViewData[key] = val
}
