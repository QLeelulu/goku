package goku

import (
    "testing"
    //"fmt"
    "github.com/sdegutis/go.assert"
)

var r1 = &Route{
    Name:    "default",
    Pattern: "/{controller}/{action}/{id}",
    Default: map[string]string{"controller": "home", "action": "index", "id": "0"},
}
var r2 = &Route{
    Name:    "nodefault",
    Pattern: "/{controller}/{action}/{id}",
}
var r3 = &Route{
    Name:       "onstraint",
    Pattern:    "/{controller}/{action}/{id}",
    Default:    map[string]string{"controller": "home", "action": "index", "id": "0"},
    Constraint: map[string]string{"id": "\\d+"},
}

func initRoute() {
    r1.Init()
    r2.Init()
    r3.Init()
}

var routerTestData = []struct {
    Route   *Route
    Url     string
    Matched bool
    Rd      *RouteData
}{
    {
        r1,
        "/",
        true,
        &RouteData{
            Controller: "home",
            Action:     "index",
            Params:     map[string]string{"id": "0"},
        },
    },
    {
        r1,
        "/admin/post/3",
        true,
        &RouteData{
            Controller: "admin",
            Action:     "post",
            Params:     map[string]string{"id": "3"},
        },
    },
    {
        r2,
        "/",
        false,
        &RouteData{},
    },
    {
        r2,
        "/home/index/a",
        true,
        &RouteData{
            Controller: "home",
            Action:     "index",
            Params:     map[string]string{"id": "a"},
        },
    },
    {
        r3,
        "/",
        true,
        &RouteData{
            Controller: "home",
            Action:     "index",
            Params:     map[string]string{"id": "0"},
        },
    },
    {
        r3,
        "/home/index/a",
        false,
        &RouteData{
            Controller: "home",
            Action:     "index",
            Params:     map[string]string{"id": "1"},
        },
    },
    {
        r3,
        "/home/index/3",
        true,
        &RouteData{
            Controller: "home",
            Action:     "index",
            Params:     map[string]string{"id": "3"},
        },
    },
}

func TestRouteInit(t *testing.T) {
    defer func() {
        if x := recover(); x != nil {
            t.Errorf("Must no panic, but got panic: \n\t%s", x)
        }
    }()
    router := new(Route)
    router.Name = "testRoute"
    router.Pattern = "/{controller}/{action}/{id}"
    router.Default = map[string]string{"controller": "home", "action": "index", "id": "0"}
    router.Init()

    rd, ok := router.Match("/home/index/2")
    assert.Equals(t, ok, true)
    assert.Equals(t, rd.Controller, "home")
    assert.Equals(t, rd.Action, "index")
    assert.Equals(t, rd.Params["id"], "2")
}

func TestRouteMatch(t *testing.T) {
    initRoute()

    for _, td := range routerTestData {
        rd, ok := td.Route.Match(td.Url)
        assert.Equals(t, ok, td.Matched)
        if td.Matched && ok {
            assert.Equals(t, rd.Controller, td.Rd.Controller)
            assert.Equals(t, rd.Action, td.Rd.Action)
            for k, p := range td.Rd.Params {
                assert.Equals(t, rd.Params[k], p)
            }
        }
    }
}

func TestRouteTable(t *testing.T) {
    var rt *RouteTable
    rt = &RouteTable{Routes: make([]*Route, 0, 10)}

    rt.Map(
        "post",
        "/post/{action}/{id}",
        map[string]string{"controller": "post"},
        map[string]string{"id": "\\d+"},
    )
    rt.AddRoute(r2)
    rt.AddRoute(r1)

    rd, ok := rt.Match("/p")
    if ok {
        assert.Equals(t, rd.Route.Name, "default")
    }
    rd, ok = rt.Match("/post/save")
    if ok {
        assert.Equals(t, rd.Route.Name, "default")
    }
    rd, ok = rt.Match("/post/update/2")
    if ok {
        assert.Equals(t, rd.Route.Name, "post")
    }
    rd, ok = rt.Match("/home/index/3")
    if ok {
        assert.Equals(t, rd.Route.Name, "nodefault")
    }
}
