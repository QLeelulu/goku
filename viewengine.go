package goku

import (
	"text/template"
	"io"
	"fmt"
)

// TemplateEnginer interface
type TemplateEnginer interface {
	Render(filepath string, viewData interface{}, wr io.Writer)
}

// DefaultTemplateEngine
type DefaultTemplateEngine struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
}

func (te *DefaultTemplateEngine) Render(filepath string, viewData interface{}, wr io.Writer) {
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

// ViewEnginer interface
type ViewEnginer interface {
	Render(controller string, action string, viewName string, viewData interface{}, wr io.Writer)
	LookupView(controller string, action string, viewName string) (filePath string)
}

// DefaultViewEngine
type DefaultViewEngine struct {
	RootDir        string
	TemplateEngine TemplateEnginer
	UseCache       bool
	Caches         map[string]string // controller & action & view to the real-file-path cache
}

func (ve *DefaultViewEngine) LookupView(controller string, action string, viewName string) string {
	if ve.UseCache {
		if v, ok := ve.Caches[controller+"_"+action+"_"+viewName]; ok {
			return v
		}
	}
	lookPaths := make([]string, 0)

	panic(fmt.Sprintf("DefaultViewEngine: can't find the view for {controller: %s, action: %s, view: %s}, look up paths: %s",
		controller, action, viewName, lookPaths))
	return ""
}

func (ve *DefaultViewEngine) Render(controller string, action string, viewName string, viewData interface{}, wr io.Writer) {
	viewFile := ve.LookupView(controller, action, viewName)
	ve.TemplateEngine.Render(viewFile, viewData, wr)
}
