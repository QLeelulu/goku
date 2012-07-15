// a golang web mvc framework, mostly like asp.net mvc. 
// Base Features:
//      + mvc (Lightweight model)
//      + route
//      + multi template engine and layout
//      + simple database api
//      + form validation
//      + filter for controller or action
//      + middleware
// 
// Example:
// 
// package main
//
// import (
//     "github.com/QLeelulu/goku"
//     "log"
//     "path"
//     "runtime"
// )
//
// /**
//  * Controller & Action
//  */
// var _ = goku.Controller("home").
//     Get("index", func(ctx *goku.HttpContext) goku.ActionResulter {
//     return ctx.View(nil)
// })
//
// // routes
// var routes []*goku.Route = []*goku.Route{
//     &goku.Route{
//         Name:    "default",
//         Pattern: "/{controller}/{action}/",
//         Default: map[string]string{"controller": "home", "action": "index"},
//     },
// }
//
// // server config
// var config *goku.ServerConfig = &goku.ServerConfig{Addr: ":8080"}
//
// func init() {
//     // project root dir, this code can not put to main func
//     _, filename, _, _ := runtime.Caller(1)
//     config.RootDir = path.Dir(filename)
// }
//
// func main() {
//     rt := &goku.RouteTable{Routes: routes}
//     s := goku.CreateServer(rt, nil, config)
//
//     goku.Logger().Logln("Server start on", s.Addr)
//     log.Fatal(s.ListenAndServe())
// }
//
package goku
