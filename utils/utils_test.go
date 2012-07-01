package utils

import (
    "testing"
    "regexp"
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

func BenchmarkNamedRegexpGroup(b *testing.B) {
    b.StopTimer()
    reg := regexp.MustCompile("(?P<name>\\w+)(?P<age>\\d+)-(\\d+)-(\\d+)")
    str := "QLeelulu25-888-333"
    b.StartTimer()
    for i := 0; i < b.N; i++ {
        NamedRegexpGroup(str, reg)
    }
}

