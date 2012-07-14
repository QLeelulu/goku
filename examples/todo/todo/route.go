package todo

import (
    "github.com/QLeelulu/goku"
)

// routes
var Routes []*goku.Route = []*goku.Route{
    &goku.Route{
        Name:     "static",
        IsStatic: true,
        Pattern:  "/public/(.*)",
    },
    &goku.Route{
        Name:       "edit",
        Pattern:    "/{controller}/{id}/{action}",
        Default:    map[string]string{"action": "edit"},
        Constraint: map[string]string{"id": "\\d+"},
    },
    &goku.Route{
        Name:    "default",
        Pattern: "/{controller}/{action}",
        Default: map[string]string{"controller": "todo", "action": "index"},
    },
}
