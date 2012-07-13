package todo

import (
    "github.com/qleelulu/goku"
)

// routes
var Routes []*goku.Route = []*goku.Route{
    &goku.Route{
        Name:     "static",
        IsStatic: true,
        Pattern:  "/public/(.*)",
    },
    &goku.Route{
        Name:       "onstraint",
        Pattern:    "/{controller}/{action}/{id}",
        Default:    map[string]string{"controller": "home", "action": "index", "id": "0"},
        Constraint: map[string]string{"id": "\\d+"},
    },
    &goku.Route{
        Name:    "default",
        Pattern: "/{controller}/{action}/{id}",
        Default: map[string]string{"controller": "home", "action": "index", "id": "0"},
    },
}
