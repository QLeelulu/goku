package utils

import (
    "testing"
    "regexp"
    "github.com/sdegutis/go.assert"
)

func TestNamedRegexpGroup(t *testing.T) {
    // (?P<name>re)
    reg := regexp.MustCompile("(?P<name>[A-Za-z]+)(?P<age>\\d+)-(\\d+)-(\\d+)")
    res, ok := NamedRegexpGroup("QLeelulu25-888-333", reg)
    if !ok {
        t.Errorf("no matched!")
    }
    if len(res) != 2 {
        t.Errorf("length error for %d", len(res))
    }
    if name, _ := res["name"]; name != "QLeelulu" {
        t.Errorf("group {name} for %s", name)
    }
}

func TestSnakeCasedName(t *testing.T) {
    ss := "HelloWorld"
    s := SnakeCasedName(ss)
    if s != "hello_world" {
        t.Errorf("wrong value for %s => %s", ss, s)
    }
}

type Blog struct {
    Author string
    Title  string
    Body   string
    Rating int
}

func TestStructToMap(t *testing.T) {
    b := Blog{
        Author: "steven",
        Title:  "stuff is cool",
        Body:   "yep! truly indeed.",
        Rating: 5,  // out of 10, tho, so...
    }

    m := StructToMap(b)

    assert.DeepEquals(t, m, map[string]interface{}{
        "Author": "steven",
        "Title":  "stuff is cool",
        "Body":   "yep! truly indeed.",
        "Rating": 5,
    })
}

func TestMapToStruct(t *testing.T) {
    m := map[string]interface{}{
        "Author": "steven",
        "Title":  "stuff is cool",
        "Body":   "yep! truly indeed.",
        "Rating": 5,
    }

    var b Blog
    MapToStruct(m, &b)

    assert.DeepEquals(t, b, Blog{
        Author: "steven",
        Title:  "stuff is cool",
        Body:   "yep! truly indeed.",
        Rating: 5,  // out of 10, tho, so...
    })
}

func TestStructName(t *testing.T) {
    var b Blog
    assert.Equals(t, StructName(b), "Blog")
    assert.Equals(t, StructName(&b), "Blog")

    var b2 *Blog = &b
    assert.Equals(t, StructName(&b2), "Blog")
}

func BenchmarkNamedRegexpGroup(b *testing.B) {
    b.StopTimer()
    reg := regexp.MustCompile("(?P<name>\\w+)(?P<age>\\d+)-(\\d+)-(\\d+)")
    str := "QLeelulu25-888-333"
    b.StartTimer()
    for i := 0; i < b.N; i++ {
        NamedRegexpGroup(str, reg)
    }
}
