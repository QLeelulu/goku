package form

import (
    "testing"
    //"fmt"
    "github.com/sdegutis/go.assert"
)

func TestNumberValidater(t *testing.T) {
    vr := &ValidResult{}
    opt := &FieldOption{}
    vd := &numberValidater{}
    vr = vd.Valid(" a123", opt)
    assert.Equals(t, vr.IsValid, false)
    assert.Equals(t, vr.CleanValue.(string), "")

    vr = vd.Valid(" 123", opt)
    assert.Equals(t, vr.IsValid, true)
    assert.Equals(t, vr.CleanValue.(string), "123")

    vr = vd.Valid(".123", opt)
    assert.Equals(t, vr.IsValid, false)
    assert.Equals(t, vr.CleanValue.(string), "")

    vr = vd.Valid("1.23", opt)
    assert.Equals(t, vr.IsValid, true)
    assert.Equals(t, vr.CleanValue.(string), "1.23")

    vr = vd.Valid("1.2.3", opt)
    assert.Equals(t, vr.IsValid, false)
    assert.Equals(t, vr.CleanValue.(string), "")

    ropt := &FieldOption{
        Required: true,
    }

    vr = vd.Valid(" ", opt)
    assert.Equals(t, vr.IsValid, true)
    assert.Equals(t, vr.CleanValue.(string), "")

    vr = vd.Valid(" ", ropt)
    assert.Equals(t, vr.IsValid, false)
    assert.Equals(t, vr.CleanValue.(string), "")

    tf := NewTextField("name", "nickname", true).Range(3, 5)
    tf.SetValue("lulu")
    vr = tf.Valid()
    assert.Equals(t, vr.IsValid, true)
    assert.Equals(t, vr.CleanValue.(string), "lulu")
    tf.SetValue("lu")
    vr = tf.Valid()
    assert.Equals(t, vr.IsValid, false)
    assert.Equals(t, vr.CleanValue.(string), "")
    tf.SetValue("qleelulu")
    vr = tf.Valid()
    assert.Equals(t, vr.IsValid, false)
    assert.Equals(t, vr.CleanValue.(string), "")
}

func TestForm(t *testing.T) {
    name := NewCharField("name", "名称", true).Range(3, 10).Field()
    nickName := NewCharField("nick_name", "昵称", false).Min(3).Max(20).Field()
    age := NewIntegerField("age", "年龄", true).Range(18, 50).Field()
    content := NewTextField("content", "内容", true).Min(10).Field()

    formData := map[string]string{
        "name":      "QLeelulu",
        "nick_name": "刘一刀",
        "age":       "22",
        "content":   "内容内容！！！！！!!!!!!",
    }

    form := NewForm(name, nickName, age, content)
    form.FillByMap(formData)
    isValid := form.Valid()
    sv := form.Values()
    cv := form.CleanValues()

    assert.Equals(t, isValid, true)
    assert.Equals(t, cv["name"], "QLeelulu")
    assert.Equals(t, cv["age"], 22)
    assert.Equals(t, sv["name"], "QLeelulu")

    formData2 := map[string]string{
        "name":      "lu",
        "nick_name": "刀",
        "age":       "d22",
        "content":   "内容内容！！",
    }
    form = NewForm(name, nickName, age, content)
    form.FillByMap(formData2)
    isValid = form.Valid()
    sv = form.Values()
    cv = form.CleanValues()

    assert.Equals(t, isValid, false)
    assert.Equals(t, cv["name"], "")
    assert.Equals(t, cv["age"], 0)
    assert.Equals(t, sv["name"], "lu")
    // fmt.Printf("%s\n", form.Errors())
}
