package goku

import (
    "strings"
    "fmt"
)

// type HttpMethod int8

// const (
// 	ALL_METHOD HttpMethod = 0
// 	GET        HttpMethod = 1
// 	POST       HttpMethod = 2
// 	PUT        HttpMethod = 4
// 	DELETE     HttpMethod = 8
// 	HEAD       HttpMethod = 16
// )

type ActionInfo struct {
    Name       string
    Controller *ControllerInfo
    Handler    func(ctx *HttpContext) ActionResulter
    Filters    []Filter
}

func (ai *ActionInfo) AddFilters(filters ...Filter) {
    for _, ft := range filters {
        if ft != nil {
            ai.Filters = append(ai.Filters, ft)
        }
    }
}

type ControllerInfo struct {
    Name    string
    Actions map[string]*ActionInfo
    Filters []Filter
}

func (ci *ControllerInfo) Init() *ControllerInfo {
    ci.Actions = make(map[string]*ActionInfo)
    return ci
}

// get a action
// e.g. ci.GetAction("get", "index"), 
// will found the registered action "index" for http method "get",
// if not found, will found the action "index" for all the http method
func (ci *ControllerInfo) GetAction(method string, name string) *ActionInfo {
    ai, ok := ci.Actions[strings.ToLower(method)+"_"+strings.ToLower(name)]
    if !ok {
        // get the action for all the http method
        ai, _ = ci.Actions["_"+strings.ToLower(name)]
    }
    return ai
}

func (ci *ControllerInfo) RegAction(httpMethod string, actionName string,
    handler func(ctx *HttpContext) ActionResulter) *ActionInfo {
    httpMethod = strings.ToLower(httpMethod)
    if httpMethod == "all" {
        httpMethod = ""
    }
    index := fmt.Sprintf("%s_%s", httpMethod, strings.ToLower(actionName))
    // check whether the action has registered
    _, ok := ci.Actions[index]
    if ok {
        panic(fmt.Sprintf("%s %s.%s has registered.", strings.ToUpper(httpMethod), ci.Name, actionName))
    }
    ai := &ActionInfo{
        Name:       strings.ToLower(actionName),
        Controller: ci,
        Handler:    handler,
    }
    ci.Actions[index] = ai
    return ai
}

// add filters for the controller
func (ci *ControllerInfo) AddFilters(filters ...Filter) {
    for _, ft := range filters {
        if ft != nil {
            ci.Filters = append(ci.Filters, ft)
        }
    }
}

// add filters for the controller
func (ci *ControllerInfo) AddActionFilters(httpMethod string, actionName string, filters ...Filter) {
    ai := ci.GetAction(httpMethod, actionName)
    if ai == nil {
        panic("ControllerInfo.AddActionFilters: controller \"" + ci.Name + "\" no action for \"" +
            strings.ToUpper(httpMethod) + " " + actionName + "\".")
    }
    ai.AddFilters(filters...)
}

type ControllerFactory struct {
    Controllers map[string]*ControllerInfo
}

func (cf *ControllerFactory) GetAction(httpMethod string, controller string, action string) *ActionInfo {
    c, ok := cf.Controllers[strings.ToLower(controller)]
    if !ok {
        return nil
    }
    return c.GetAction(httpMethod, action)
}

var defaultControllerFactory *ControllerFactory = &ControllerFactory{
    Controllers: make(map[string]*ControllerInfo),
}

type ControllerBuilder struct {
    controller    *ControllerInfo
    currentAction *ActionInfo
}

// @param httpMethod: if "all", will match all http method, but Priority is low
func (cb *ControllerBuilder) Action(httpMethod string, actionName string,
    handler func(ctx *HttpContext) ActionResulter) *ControllerBuilder {

    cb.currentAction = cb.controller.RegAction(httpMethod, actionName, handler)
    return cb
}

// reg http "get" method action
func (cb *ControllerBuilder) Get(actionName string,
    handler func(ctx *HttpContext) ActionResulter) *ControllerBuilder {
    return cb.Action("get", actionName, handler)
}

// reg http "post" method action
func (cb *ControllerBuilder) Post(actionName string,
    handler func(ctx *HttpContext) ActionResulter) *ControllerBuilder {

    return cb.Action("post", actionName, handler)
}

// reg http "put" method action
func (cb *ControllerBuilder) Put(httpMethod string, actionName string,
    handler func(ctx *HttpContext) ActionResulter) *ControllerBuilder {

    return cb.Action("put", actionName, handler)
}

// reg http "delete" method action
func (cb *ControllerBuilder) Delete(httpMethod string, actionName string,
    handler func(ctx *HttpContext) ActionResulter) *ControllerBuilder {

    return cb.Action("delete", actionName, handler)
}

func (cb *ControllerBuilder) Filters(filters ...Filter) *ControllerBuilder {
    if cb.currentAction != nil {
        cb.currentAction.AddFilters(filters...)
    } else {
        cb.controller.AddFilters(filters...)
    }
    return cb
}

func Controller(name string) *ControllerBuilder {
    name = strings.ToLower(name)
    c, ok := defaultControllerFactory.Controllers[name]
    if !ok {
        c = &ControllerInfo{
            Name: name,
        }
        c.Init()
        // add to index
        defaultControllerFactory.Controllers[name] = c
    }
    cb := &ControllerBuilder{
        controller: c,
    }
    return cb
}
