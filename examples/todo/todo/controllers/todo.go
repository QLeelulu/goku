package controllers

import (
    // "fmt"
    "github.com/QLeelulu/goku"
    "github.com/QLeelulu/goku/examples/todo/todo/models"
    "github.com/QLeelulu/goku/form"
    "strconv"
    "time"
)

/**
 * form
 */
func createTodoForm() *form.Form {
    // defined the field
    id := form.NewIntegerField("id", "Id", false).Range(1, 10).
        // {0} will replace with the min range value,
        // {1} will replace with the max range value,
        Error("range", "值必须在{0}到{1}之间").Field()
    // title
    title := form.NewTextField("title", "待办事项", true).Min(8).Max(200).
        Error("required", "必须填写事项内容").
        Error("range", "字数必须在{0}到{1}之间").Field()

    // add the fields to a form
    form := form.NewForm(id, title)
    return form
}

/**
 * todo controller
 */
var _ = goku.Controller("todo").

    // todo.index action
    Get("index", func(ctx *goku.HttpContext) goku.ActionResulter {

    todos, err := models.GetTodoLists()
    if err != nil {
        return ctx.Error(err)
    }
    // you can pass a struct to ViewModel
    // then you can use it in template
    // like this: {{ .Model.xxx }}
    return ctx.View(todos)
}).
    /**
     * todo.new action
     */
    Post("new", func(ctx *goku.HttpContext) goku.ActionResulter {

    f := createTodoForm()
    f.FillByRequest(ctx.Request)

    errorMsgs := ""
    // valid the form values
    if f.Valid() {
        // after valid, we can get the clean values
        m := f.CleanValues()
        delete(m, "id")
        m["post_date"] = time.Now()
        // save the value to db, see models/todo.go
        _, err := models.SaveTodo(m)
        if err == nil {
            return ctx.Redirect("/")
        }
        errorMsgs = "Database error"
        goku.Logger().Errorln(err)
    } else {
        // if not valid
        // we can get the valid errors
        errs := f.Errors()
        for _, v := range errs {
            if errorMsgs != "" {
                errorMsgs += ","
            }
            errorMsgs += v[0] + ": " + v[1]
        }
    }
    // you can add any val to ViewData
    // then you can use it in template
    // like this: {{ .Data.errorMsg }}
    ctx.ViewData["errorMsg"] = errorMsgs
    return ctx.Render("error", nil)
}).
    /**
     * todo.edit action for get
     */
    Get("edit", func(ctx *goku.HttpContext) goku.ActionResulter {

    id, err := strconv.Atoi(ctx.RouteData.Params["id"])
    if err == nil {
        var todo models.Todo
        todo, err = models.GetTodo(id)
        if err == nil {
            return ctx.View(todo)
        }
    }
    ctx.ViewData["errorMsg"] = err.Error()
    return ctx.Render("error", nil)
}).
    /**
     * todo.edit action fot post
     */
    Post("edit", func(ctx *goku.HttpContext) goku.ActionResulter {

    f := createTodoForm()
    f.FillByRequest(ctx.Request)

    errorMsgs := ""
    // valid the form values
    if f.Valid() {
        // after valid, we can get the clean values
        m := f.CleanValues()
        id := m["id"].(int)
        delete(m, "id")
        // update the value to db, see models/todo.go
        _, err := models.UpdateTodo(id, m)
        if err == nil {
            return ctx.Redirect("/")
        }
        errorMsgs = "Database error"
        goku.Logger().Errorln(err)
    } else {
        // if not valid
        // we can get the valid errors
        errs := f.Errors()
        for _, v := range errs {
            if errorMsgs != "" {
                errorMsgs += ","
            }
            errorMsgs += v[0] + ": " + v[1]
        }
    }
    ctx.ViewData["errorMsg"] = errorMsgs
    return ctx.Render("error", nil)
}).
    /**
     * todo.finish action
     */
    Get("finish", func(ctx *goku.HttpContext) goku.ActionResulter {

    id, err := strconv.Atoi(ctx.RouteData.Params["id"])
    status := ctx.Request.FormValue("status")

    errorMsgs := ""
    if id > 0 && (status == "yes" || status == "no") {
        db := models.GetDB()
        defer db.Close()
        fsh := 0
        if status == "yes" {
            fsh = 1
        }
        updateVals := map[string]interface{}{"finished": fsh}
        _, err = db.Update("todo", updateVals, "id=?", id)
        if err == nil {
            return ctx.Redirect("/")
        }
        errorMsgs = "Database error"
        goku.Logger().Errorln(err)
    } else {
        errorMsgs = "错误的请求"
    }
    ctx.ViewData["errorMsg"] = errorMsgs
    return ctx.Render("error", nil)
}).
    /**
     * todo.delete action
     */
    Get("delete", func(ctx *goku.HttpContext) goku.ActionResulter {

    id, err := strconv.Atoi(ctx.RouteData.Params["id"])

    errorMsg := ""
    if id > 0 {
        _, err = models.DeleteTodo(id)
        if err == nil {
            return ctx.Redirect("/")
        }
        errorMsg = "Database error"
        goku.Logger().Errorln(err)
    } else {
        errorMsg = "错误的请求"
    }
    ctx.ViewData["errorMsg"] = errorMsg
    return ctx.Render("error", nil)
})
