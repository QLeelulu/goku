package utils

import (
    "regexp"
    "reflect"
    // "fmt"
    "os"
)

// match regexp with string, and return a named group map
// Example:
//   regexp: "(?P<name>[A-Za-z]+)-(?P<age>\\d+)"
//   string: "CGC-30"
//   return: map[string]string{ "name":"CGC", "age":"30" }
func NamedRegexpGroup(str string, reg *regexp.Regexp) (ng map[string]string, matched bool) {
    rst := reg.FindStringSubmatch(str)
    //fmt.Printf("%s => %s => %s\n\n", reg, str, rst)
    if len(rst) < 1 {
        return
    }
    ng = make(map[string]string)
    lenRst := len(rst)
    sn := reg.SubexpNames()
    for k, v := range sn {
        // SubexpNames contain the none named group,
        // so must filter v == ""
        if k == 0 || v == "" {
            continue
        }
        if k+1 > lenRst {
            break
        }
        ng[v] = rst[k]
    }
    matched = true
    return
}

// check if that the file exists
func FileExists(path_ string) (bool, error) {
    _, err := os.Stat(path_)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}

// convert like this: "HelloWorld" to "hello_world"
func SnakeCasedName(name string) string {
    newstr := make([]rune, 0)
    firstTime := true

    for _, chr := range name {
        if isUpper := 'A' <= chr && chr <= 'Z'; isUpper {
            if firstTime == true {
                firstTime = false
            } else {
                newstr = append(newstr, '_')
            }
            chr -= ('A' - 'a')
        }
        newstr = append(newstr, chr)
    }

    return string(newstr)
}

// map & struct convert is from https://github.com/sdegutis/go.mapstruct

// convert map to struct
func MapToStruct(m map[string]interface{}, s interface{}) {
    v := reflect.Indirect(reflect.ValueOf(s))

    for i := 0; i < v.NumField(); i++ {
        key := v.Type().Field(i).Name
        v.Field(i).Set(reflect.ValueOf(m[key]))
    }
}

// convert struct to map
func StructToMap(s interface{}) map[string]interface{} {
    m := make(map[string]interface{})

    v := reflect.ValueOf(s)

    for i := 0; i < v.NumField(); i++ {
        key := v.Type().Field(i).Name
        val := v.Field(i).Interface()

        m[key] = val
    }
    return m
}

// get the Struct's name
func StructName(s interface{}) string {
    v := reflect.TypeOf(s)
    for v.Kind() == reflect.Ptr {
        v = v.Elem()
    }
    return v.Name()
}
