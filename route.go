package goku

// all the thing about route
// may be change route to interface ?

import (
    "fmt"
    "regexp"
    //"path"
    "github.com/QLeelulu/goku/utils"
)

var (
    regPathParse *regexp.Regexp = regexp.MustCompile("/?\\{[\\w\\-_]+\\}") // matched like this: /{controller}
)

// Route config
//      var rt = &Route {
//          Name: "default",
//          Pattern: "/{controller}/{action}/{id}",
//          Default: map[string]string { "controller": "home", "action": "index", "id": "0", },
//          Constraint: map[string]string { "id": "\\d+" }
//      }
// and then, must init the router
//      rt.Init()
// and then, you can use it
//      rt.Match("/home/index")
//
type Route struct {
    Name       string            // the router name
    Pattern    string            // url pattern config, eg. /{controller}/{action}/{id}
    Default    map[string]string // default value for Pattern
    Constraint map[string]string // constraint for Pattern, value is regexp str
    IsStatic   bool              // whether the route is for static file

    rePath *regexp.Regexp
    inited bool
}

func (router *Route) Init() {
    if router.inited {
        return
    }
    if router.Name == "" {
        panic("Route: Name must be set")
    }
    if router.Pattern == "" {
        panic("Route: Pattern must be set")
    }
    if router.Default == nil {
        router.Default = make(map[string]string)
    }
    if router.Constraint == nil {
        router.Constraint = make(map[string]string)
    }

    //  /{controller}/{action}/{id} 
    //      => /(?P<controller>[^\?#/]+)/(?P<action>[^\?#/]+)/(?P<id>[^\?#/]+)?
    r := regPathParse.ReplaceAllStringFunc(router.Pattern, func(s string) string {
        slash := ""
        if s[0] == '/' {
            slash = "/"
            s = s[1:]
        }
        name, reg, need := s[1:len(s)-1], "", ""
        if v, ok := router.Constraint[name]; ok {
            reg = v
        } else {
            if name == "controller" || name == "action" {
                reg = "[^\\.\\?#/]+"
            } else {
                reg = "[^\\?#/]+"
            }
        }
        // if default value exist, it's options
        if _, ok := router.Default[name]; ok {
            need = "?"
            if slash != "" {
                //  / => /?
                slash = slash + "?"
            }
        }
        //(?P<name>re)
        return fmt.Sprintf("%s(?P<%s>%s)%s", slash, name, reg, need)
    })
    if r != "" && r[len(r)-1] == '/' {
        r = r + "?"
    }
    router.rePath = regexp.MustCompile("^" + r + "$")
    router.inited = true
}

func (router *Route) Match(url string) (rd *RouteData, matched bool) {
    if !router.inited {
        router.Init()
    }
    if router.IsStatic {
        rd, matched = router.matchStatic(url)
        if matched {
            return
        }
    }
    md, ok := utils.NamedRegexpGroup(url, router.rePath)
    if !ok {
        return
    }
    if md["controller"] == "" && router.Default["controller"] == "" {
        return
    }
    if md["action"] == "" && router.Default["action"] == "" {
        return
    }

    rd = new(RouteData)
    rd.Url = url
    rd.Route = router

    for k, v := range router.Default {
        if md[k] == "" {
            md[k] = v
        }
    }

    rd.Controller = md["controller"]
    rd.Action = md["action"]
    delete(md, "controller")
    delete(md, "action")
    rd.Params = md

    matched = true
    return
}

// static file route match
// if has group ,return group 1, else return the url
// e.g.
//		pattern: /static/.*  , url: /static/logo.gif , static path: /static/logo.gif
//		pattern: /static/(.*)  , url: /static/logo.gif , static path: logo.gif
func (route *Route) matchStatic(url string) (rd *RouteData, matched bool) {
    rst := route.rePath.FindStringSubmatch(url)
    lrst := len(rst)
    if lrst < 1 {
        return
    }
    rd = new(RouteData)
    rd.Url = url
    rd.Route = route
    if lrst > 1 {
        rd.FilePath = rst[1]
    } else {
        rd.FilePath = rst[0]
    }
    matched = true
    return
}

type RouteData struct {
    Url        string
    Route      *Route // is this field need ?
    Controller string
    Action     string
    Params     map[string]string
    FilePath   string // if is a static file route, this will be set
}

func (rd *RouteData) Get(name string) (val string, ok bool) {
    val, ok = rd.Params[name]
    return
}

type RouteTable struct {
    Routes []*Route
}

func (rt *RouteTable) Match(url string) (rd *RouteData, matched bool) {
    if url == "" {
        return
    }
    for _, router := range rt.Routes {
        if _rd, ok := router.Match(url); ok {
            rd = _rd
            matched = true
            break
        }
    }
    return
}

func (rt *RouteTable) AddRoute(route *Route) {
    route.Init()
    rt.Routes = append(rt.Routes, route)
}

// add a new route
// params:
//	+ name: route name
//  + url:  url pattern
//  + default: map[string]string, default value for url pattern
//  + constraint: map[string]string, constraint for url pattern
func (rt *RouteTable) Map(name string, url string, args ...interface{}) {
    var defaultData map[string]string
    var constraint map[string]string
    if len(args) > 0 {
        defaultData = args[0].(map[string]string)
    }
    if len(args) > 1 {
        constraint = args[1].(map[string]string)
    }

    route := &Route{
        Name:       name,
        Pattern:    url,
        Default:    defaultData,
        Constraint: constraint,
    }
    route.Init()
    rt.Routes = append(rt.Routes, route)
}

// static file route match
// if has group ,return group 1, else return the url
// e.g.
//		pattern: /static/.*  , url: /static/logo.gif , static path: /static/logo.gif
//		pattern: /static/(.*)  , url: /static/logo.gif , static path: logo.gif
func (rt *RouteTable) Static(name string, pattern string) {
    route := &Route{
        Name:     name,
        IsStatic: true,
        Pattern:  pattern,
    }
    route.Init()
    rt.Routes = append(rt.Routes, route)
}
