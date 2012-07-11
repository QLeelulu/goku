package form

import (
    "strings"
    "strconv"
    "regexp"
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
    Max      int
    Min      int
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
    if opt.Required && source == "" {
        vr.ErrorMsg = MSG_REQUIRED
    }
    return source, vr
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
            if opt.ErrorMsg != "" {
                vr.ErrorMsg = opt.ErrorMsg
            } else {
                vr.ErrorMsg = "not a valid number"
            }
            return vr
        }
    }
    vr.IsValid = true
    vr.CleanValue = source
    return vr
}

// check if the source value is match the givent regexp
type regexpValidater struct {
    baseValidater

    Regexp *regexp.Regexp
}

func (rv *regexpValidater) Valid(source string, opt *FieldOption) (vr *ValidResult) {
    source, vr = rv.baseValidater.Valid(source, opt)
    if vr.ErrorMsg != "" {
        return vr
    }
    ok := rv.Regexp.MatchString(source)
    if ok {
        vr.IsValid = true
        vr.CleanValue = source
    } else {
        if opt.ErrorMsg != "" {
            vr.ErrorMsg = opt.ErrorMsg
        } else {
            vr.ErrorMsg = "not match"
        }
    }
    return vr
}

// check if the source value is int value
type intValidater struct {
    baseValidater
}

func (nv *intValidater) Valid(source string, opt *FieldOption) (vr *ValidResult) {
    source, vr = nv.baseValidater.Valid(source, opt)
    if vr.ErrorMsg != "" {
        return vr
    }
    val, err := strconv.Atoi(source)
    vr.CleanValue = val
    if err == nil {
        vr.IsValid = true
    } else {
        if opt.ErrorMsg != "" {
            vr.ErrorMsg = opt.ErrorMsg
        } else {
            vr.ErrorMsg = "not a int value"
        }
    }
    return vr
}

type stringValidater struct {
    baseValidater
}

func (nv *stringValidater) Valid(source string, opt *FieldOption) (vr *ValidResult) {
    source, vr = nv.baseValidater.Valid(source, opt)
    if vr.ErrorMsg != "" {
        return vr
    }
    if opt.Range[0] > 0 && len(source) < opt.Range[0] {
        if opt.Range[1] > 0 {
            vr.ErrorMsg = strings.Replace(MSG_RANGE, "{0}", strconv.Itoa(opt.Range[0]), -1)
            vr.ErrorMsg = strings.Replace(vr.ErrorMsg, "{1}", strconv.Itoa(opt.Range[1]), -1)
        } else {
            vr.ErrorMsg = strings.Replace(MSG_MIN_LENGTH, "{0}", strconv.Itoa(opt.Range[0]), -1)
        }
        return vr
    }
    if opt.Range[1] > 0 && len(source) > opt.Range[1] {
        if opt.Range[0] > 0 {
            vr.ErrorMsg = strings.Replace(MSG_RANGE, "{0}", strconv.Itoa(opt.Range[0]), -1)
            vr.ErrorMsg = strings.Replace(vr.ErrorMsg, "{1}", strconv.Itoa(opt.Range[1]), -1)
        } else {
            vr.ErrorMsg = strings.Replace(MSG_MAX_LENGTH, "{0}", strconv.Itoa(opt.Range[1]), -1)
        }
        return vr
    }
    vr.IsValid = true
    vr.CleanValue = source

    return vr
}
