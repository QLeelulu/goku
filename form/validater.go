package form

import (
    "regexp"
    "strconv"
    "strings"
    "unicode/utf8"
)

const (
    MSG_REQUIRED     = "this field is required"
    MSG_MAX_LENGTH   = "length must less than {0}"
    MSG_MIN_LENGTH   = "length must more than {0}"
    MST_RANGE_LENGTH = "length must between {0} and {1}"
    MSG_MAX          = "value must small than {0}"
    MSG_MIN          = "value must large than {0}"
    MSG_RANGE        = "value must between {0} and {1}"
    MSG_INVALID      = "not a valid value"
)

type ValidResult struct {
    IsValid    bool
    ErrorMsg   string
    CleanValue interface{}
}

type ValidOption struct {
    Required bool
    NotTrim  bool // not trim the Whitespace
    // Max      int
    // Min      int
    Range    [2]int
    ErrorMsg string
}

type Validater interface {
    Valid(source string, option *FieldOption) *ValidResult
}

// func (f *Field) IsValid() *ValidResult {
//     return f.Validater.Valid(f.Source, f.Option)
// }

type baseValidater struct {
}

func (nv *baseValidater) Valid(source string, opt *FieldOption) (string, *ValidResult) {
    vr := &ValidResult{
        CleanValue: "",
    }
    if opt == nil {
        opt = &FieldOption{}
    }

    if !opt.NotTrim {
        source = strings.TrimSpace(source)
    }
    if source == "" {
        if opt.Required {
            if msg, ok := opt.ErrorMsgs["required"]; ok {
                vr.ErrorMsg = msg
            } else {
                vr.ErrorMsg = MSG_REQUIRED
            }
        } else {
            vr.IsValid = true
        }
    }
    return source, vr
}

// if key in m, return m[key], else return defaultVal
func getOrDefault(m map[string]string, key, defaultVal string) string {
    msg, ok := m[key]
    if ok && msg != "" {
        return msg
    }
    return defaultVal
}

// only check if the source value is number format
// but still return string
// "0123" and "1.23" are valid
type numberValidater struct {
    baseValidater
}

func (nv *numberValidater) Valid(source string, opt *FieldOption) (vr *ValidResult) {
    source, vr = nv.baseValidater.Valid(source, opt)
    if vr.ErrorMsg != "" {
        return vr
    }

    dotCount := 0
    for i, c := range source {
        if c < 48 || c > 57 {
            if c == 46 && i != 0 && dotCount == 0 {
                dotCount++
                continue
            }
            vr.ErrorMsg = getOrDefault(opt.ErrorMsgs, "invalid", "not a valid number")
            return vr
        }
    }
    vr.IsValid = true
    vr.CleanValue = source
    return vr
}

// check if the source value is match the givent regexp
// error msg keys:
//      required
//      invalid
type regexpValidater struct {
    baseValidater

    Regexp *regexp.Regexp
}

func (rv *regexpValidater) Valid(source string, opt *FieldOption) (vr *ValidResult) {
    source, vr = rv.baseValidater.Valid(source, opt)
    if vr.IsValid || vr.ErrorMsg != "" {
        return vr
    }
    ok := rv.Regexp.MatchString(source)
    if ok {
        vr.IsValid = true
        vr.CleanValue = source
    } else {
        vr.ErrorMsg = getOrDefault(opt.ErrorMsgs, "invalid", "not match")
    }
    return vr
}

// check if the source value is int value
type intValidater struct {
    baseValidater
}

func (nv *intValidater) Valid(source string, opt *FieldOption) (vr *ValidResult) {
    source, vr = nv.baseValidater.Valid(source, opt)
    if vr.IsValid || vr.ErrorMsg != "" {
        if vr.IsValid {
            vr.CleanValue = 0
        }
        return vr
    }
    val, err := strconv.ParseInt(source, 10, 64)
    vr.CleanValue = val
    if err == nil {
        if opt.Range[0] > 0 && val < int64(opt.Range[0]) {
            if opt.Range[1] > 0 {
                vr.ErrorMsg = strings.Replace(
                    getOrDefault(opt.ErrorMsgs, "range", MSG_RANGE),
                    "{0}", strconv.Itoa(opt.Range[0]), -1)
                vr.ErrorMsg = strings.Replace(vr.ErrorMsg, "{1}", strconv.Itoa(opt.Range[1]), -1)
            } else {
                vr.ErrorMsg = strings.Replace(
                    getOrDefault(opt.ErrorMsgs, "min", MSG_MIN), "{0}", strconv.Itoa(opt.Range[0]), -1)
            }
            return vr
        }
        if opt.Range[1] > 0 && val > int64(opt.Range[1]) {
            if opt.Range[0] > 0 {
                vr.ErrorMsg = strings.Replace(
                    getOrDefault(opt.ErrorMsgs, "range", MSG_RANGE), "{0}", strconv.Itoa(opt.Range[0]), -1)
                vr.ErrorMsg = strings.Replace(vr.ErrorMsg, "{1}", strconv.Itoa(opt.Range[1]), -1)
            } else {
                vr.ErrorMsg = strings.Replace(
                    getOrDefault(opt.ErrorMsgs, "max", MSG_MAX), "{0}", strconv.Itoa(opt.Range[1]), -1)
            }
            return vr
        }
        vr.IsValid = true
    } else {
        vr.ErrorMsg = getOrDefault(opt.ErrorMsgs, "invalid", "not a int valid")
    }
    return vr
}

// error msg keys:
//      required
//      range
//      min
//      max
type stringValidater struct {
    baseValidater
}

func (nv *stringValidater) Valid(source string, opt *FieldOption) (vr *ValidResult) {
    source, vr = nv.baseValidater.Valid(source, opt)
    if vr.IsValid || vr.ErrorMsg != "" {
        return vr
    }
    if opt.Range[0] > 0 && utf8.RuneCountInString(source) < opt.Range[0] {
        if opt.Range[1] > 0 {
            vr.ErrorMsg = strings.Replace(
                getOrDefault(opt.ErrorMsgs, "range", MST_RANGE_LENGTH),
                "{0}", strconv.Itoa(opt.Range[0]), -1)
            vr.ErrorMsg = strings.Replace(vr.ErrorMsg, "{1}", strconv.Itoa(opt.Range[1]), -1)
        } else {
            vr.ErrorMsg = strings.Replace(
                getOrDefault(opt.ErrorMsgs, "min", MSG_MIN_LENGTH), "{0}", strconv.Itoa(opt.Range[0]), -1)
        }
        return vr
    }
    if opt.Range[1] > 0 && utf8.RuneCountInString(source) > opt.Range[1] {
        if opt.Range[0] > 0 {
            vr.ErrorMsg = strings.Replace(
                getOrDefault(opt.ErrorMsgs, "range", MST_RANGE_LENGTH), "{0}", strconv.Itoa(opt.Range[0]), -1)
            vr.ErrorMsg = strings.Replace(vr.ErrorMsg, "{1}", strconv.Itoa(opt.Range[1]), -1)
        } else {
            vr.ErrorMsg = strings.Replace(
                getOrDefault(opt.ErrorMsgs, "max", MSG_MAX_LENGTH), "{0}", strconv.Itoa(opt.Range[1]), -1)
        }
        return vr
    }
    vr.IsValid = true
    vr.CleanValue = source

    return vr
}
