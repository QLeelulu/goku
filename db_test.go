package goku

import (
    "testing"
    // "fmt"
    "time"
    "github.com/sdegutis/go.assert"
    _ "github.com/ziutek/mymysql/godrv"
)

type Todo struct {
    Id       int
    Title    string
    Finished bool
    PostDate time.Time
}

var db *MysqlDB

func TestOpenDB(t *testing.T) {
    var err error
    db, err = OpenMysql("mymysql", "tcp:localhost:3306*todo/lulu/123456")
    if err != nil {
        t.Error("can not open db")
    }
}

func TestMysqlDB(t *testing.T) {
    r, err := db.Select("todo", SqlQueryInfo{})
    assert.Equals(t, err, nil)

    var id, finished int
    var title string
    var post_date time.Time
    for r.Next() {
        err2 := r.Scan(&id, &title, &finished, &post_date)
        // fmt.Printf("a => %s => %s\n", id, err2)
        assert.Equals(t, err2, nil)
        assert.Equals(t, id > 0, true)
    }
    assert.Equals(t, err, nil)
}

func TestGetStruct(t *testing.T) {
    var todo *Todo
    err := db.GetStruct(todo, "id=?", 5)
    assert.StringContains(t, err.Error(), "struct can not be nil")

    err = db.GetStruct(Todo{}, "id=?", 5)
    assert.StringContains(t, err.Error(), "struct must be a pointer")

    todo = &Todo{}
    err = db.GetStruct(todo, "id=?", 5)
    assert.Equals(t, err, nil)
}

func TestGetStructs(t *testing.T) {
    qi := SqlQueryInfo{}
    todos, err := db.GetStructs(&Todo{}, qi)
    assert.Equals(t, err, nil)

    for _, todo_ := range todos {
        todo := todo_.(*Todo)
        assert.Equals(t, todo.Id > 0, true)
    }
}
