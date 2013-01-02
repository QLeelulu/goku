---
layout: post
title: "Form Validation"
description: ""
category: doc
tags: []
---
{% include JB/setup %}

you can create a form, to valid the user's input, and get the clean value.

{% highlight go %}
import "github.com/QLeelulu/goku/form"

func CreateCommentForm() *goku.Form {
    name := NewCharField("name", "Name", true).Range(3, 10).Field()
    nickName := NewCharField("nick_name", "Nick Name", false).Min(3).Max(20).Field()
    age := NewIntegerField("age", "Age", true).Range(18, 50).Field()
    content := NewTextField("content", "Content", true).Min(10).Field()

    form := NewForm(name, nickName, age, content)
    return form
}
{% endhighlight %}

and then you can use this form like this:

{% highlight go %}
f := CreateCommentForm()
f.FillByRequest(ctx.Request)

if f.Valid() {
    // after valid, we can get the clean values
    m := f.CleanValues()
    // and now you can save m to database
} else {
    // if not valid
    // we can get the valid errors
    errs := f.Errors()
}
{% endhighlight %}

checkout [form_test.go](https://github.com/QLeelulu/goku/blob/master/form/form_test.go)


