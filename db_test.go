package goku

import (
    "testing"
    // "fmt"
    "github.com/sdegutis/go.assert"
    _ "github.com/ziutek/mymysql/godrv"
    "time"
)

type TestBlog struct {
    Id       int
    Title    string
    Content  string
    CreateAt time.Time
}

var db *MysqlDB

/**
CREATE TABLE `test_blog` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(30) NOT NULL DEFAULT '',
  `content` text NOT NULL,
  `create_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
*/

func TestOpenMysqlDB(t *testing.T) {
    var err error
    db, err = OpenMysql("mymysql", "tcp:localhost:3306*test_db/lulu/123456")
    // db.Debug = true

    if err != nil {
        t.Error("can not open db")
    }
    _, err = db.Query("select 1")
    if err != nil {
        t.Error(err.Error())
    }
}

func TestMysqlDB(t *testing.T) {
    r, err := db.Select("test_blog", SqlQueryInfo{
        Fields: "id, title, content",
        Where:  "id>?",
        Params: []interface{}{0},
        Limit:  10,
        Offset: 0,
        Group:  "",
        Order:  "id desc",
    })
    if err != nil {
        t.Fatalf("select got err: ", err.Error())
    }

    var id int
    for r.Next() {
        err2 := r.Scan(&id, nil, nil)
        assert.Equals(t, err2, nil)
        assert.Equals(t, id, 1)
    }
}

func TestMysqlDBInsert(t *testing.T) {
    vals := map[string]interface{}{
        "title": "golang",
        "content": "Go is an open source programming environment that " +
            "makes it easy to build simple, reliable, and efficient software.",
        "create_at": time.Now(),
    }

    r, err := db.Insert("test_blog", vals)
    assert.Equals(t, err, nil)
    if r == nil {
        t.Fatalf("insert result must not be nil")
    }
    var i int64
    i, err = r.RowsAffected()
    assert.Equals(t, err, nil)
    assert.Equals(t, i, int64(1))
    i, err = r.LastInsertId()
    assert.Equals(t, err, nil)
    assert.Equals(t, i > 0, true)

    blog := TestBlog{
        Title:    "goku",
        Content:  "a mvc framework",
        CreateAt: time.Now(),
    }
    r, err = db.InsertStruct(&blog)
    assert.Equals(t, err, nil)
    if r == nil {
        t.Fatalf("insert result must not be nil")
    }

    i, err = r.RowsAffected()
    i, err = r.RowsAffected()
    assert.Equals(t, err, nil)
    assert.Equals(t, i, int64(1))
    i, err = r.LastInsertId()
    assert.Equals(t, err, nil)
    assert.Equals(t, blog.Id > 0, true)
}

func TestGetStruct(t *testing.T) {
    var blog *TestBlog
    var err error
    err = db.GetStruct(blog, "id=?", 1)
    // assert.Equals(t, todo, nil)
    assert.NotEquals(t, err, nil)
    if err != nil {
        assert.StringContains(t, err.Error(), "struct can not be nil")
    }

    err = db.GetStruct(TestBlog{}, "id=?", 1)
    assert.StringContains(t, err.Error(), "struct must be a pointer")

    blog = &TestBlog{}
    err = db.GetStruct(blog, "id>0")
    assert.Equals(t, err, nil)
    assert.Equals(t, blog.Id > 0, true)
}

func TestGetStructs(t *testing.T) {
    qi := SqlQueryInfo{}
    var blogs []TestBlog
    err := db.GetStructs(&blogs, qi)
    assert.Equals(t, err, nil)

    for _, blog_ := range blogs {
        blog := blog_
        assert.Equals(t, blog.Id > 0, true)
    }
}

func TestMysqlUpdate(t *testing.T) {
    var blog *TestBlog
    var err error

    blog = &TestBlog{}
    err = db.GetStruct(blog, "id>0")
    assert.Equals(t, err, nil)

    vals := map[string]interface{}{
        "title": "js",
    }
    r, err2 := db.Update("test_blog", vals, "id=?", blog.Id)
    if err2 != nil {
        t.Fatalf("update got err: ", err2.Error())
    }
    assert.Equals(t, err2, nil)
    var i int64
    i, err = r.RowsAffected()
    assert.Equals(t, err, nil)
    assert.Equals(t, i, int64(1))

    err = db.GetStruct(blog, "id=?", blog.Id)
    assert.Equals(t, err, nil)
    assert.Equals(t, blog.Title, "js")
}

func TestMysqlDBCount(t *testing.T) {
    r, err := db.Count("test_blog", "")
    assert.Equals(t, err, nil)
    assert.Equals(t, r, int64(2))
}

func TestMysqlDBDelete(t *testing.T) {
    r, err := db.Delete("test_blog", "0=0")
    assert.Equals(t, err, nil)
    if r == nil {
        t.Fatalf("delete result must not be nil")
    }

    var i int64
    i, err = r.RowsAffected()
    assert.Equals(t, err, nil)
    assert.Equals(t, i > 0, true)
}
