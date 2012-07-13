package todo

import (
    "../../../../goku"
    "time"
    "runtime"
    "path"
)

var (
    DATABASE_Driver string = "mysql"
    // "user:password@/dbname?charset=utf8"
    DATABASE_DSN string = "lulu:123456@/todo?charset=utf8&keepalive=1"
)

var Config *goku.ServerConfig = &goku.ServerConfig{
    Addr:           ":8080",
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
    //RootDir:        _, filename, _, _ := runtime.Caller(1),
    StaticPath: "static", // static content 
    ViewPath:   "views",
    Debug:      true,
}

func init() {
    // WTF, i just want to set the RootDir as current dir.
    _, filename, _, _ := runtime.Caller(1)
    Config.RootDir = path.Dir(filename)
}
