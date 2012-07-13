package goku

import (
    "fmt"
    "database/sql"
    "strings"
    "errors"
    "reflect"
    "github.com/qleelulu/goku/utils"
)

var zeroVal reflect.Value

type SQLLiteral string

type SqlQueryInfo struct {
    Fields string
    Where  string
    Params []interface{}
    Limit  int
    Offset int
    Group  string
    Order  string
}

// base db
type DB struct {
    sql.DB
}

func (db *DB) whereSql(where string) string {
    if where == "" {
        return ""
    }
    return " WHERE " + where
}

func (db *DB) orderSql(order string) string {
    if order == "" {
        return ""
    }
    return " ORDER BY " + order
}

func (db *DB) groupSql(group string) string {
    if group == "" {
        return ""
    }
    return " GROUP BY " + group
}

func (db *DB) limitSql(limit int, offset int) string {
    var r string
    if limit > 0 {
        r = fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
    }
    return r
}

// select from db.table with qi
// Example:
//      qi := &SqlQueryInfo{
//              Fields: "*",
//              Where: "id > ?",
//              Params: []interface{}{ 3 }
//              Limit: 10,
//              Offset: 0,
//              Group: "age",
//              Order: "id desc",
//      }
//      rows, err := db.Select("blog", qi)
func (db *DB) Select(table string, qi SqlQueryInfo) (*sql.Rows, error) {
    if qi.Fields == "" {
        qi.Fields = "*"
    }
    query := fmt.Sprintf("SELECT %s FROM `%s` %s %s %s %s;",
        qi.Fields,
        table,
        db.whereSql(qi.Where),
        db.groupSql(qi.Group),
        db.orderSql(qi.Order),
        db.limitSql(qi.Limit, qi.Offset),
    )
    return db.Query(query, qi.Params...)
}

// insert into table with values from vals
// Example:
//      data := map[string]interface{}{
//          "title": "hello golang",
//          "content": "just wonderful",
//      }
//      rerult, err := db.Insert("blog", data)
//      id, err := result.LastInsertId()
func (db *DB) Insert(table string, vals map[string]interface{}) (result sql.Result, err error) {
    l := len(vals)
    if vals == nil || l < 1 {
        return
    }
    fields := make([]string, 0, l)
    values := make([]string, 0, l)
    params := make([]interface{}, 0, l)
    for k, v := range vals {
        fields = append(fields, k)
        switch v.(type) {
        case SQLLiteral:
            values = append(values, fmt.Sprintf("%s", v))
        default:
            values = append(values, "?")
            params = append(params, v)
        }
    }
    query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);",
        table,
        strings.Join(fields, ", "),
        strings.Join(values, ", "))

    result, err = db.Exec(query, params...)
    return
}

// insert struct to database
// if i is pointer to struct and has a int type field named "Id"
// the field "Id" will set to the last insert id if has LastInsertId
// 
// field mapping rule is: HelloWorld => hello_world
// mean that struct's field "HelloWorld" in database table's field is "hello_world"
// table name mapping use the same rule as field
func (db *DB) InsertStruct(i interface{}) (sql.Result, error) {
    m := utils.StructToSnakeKeyMap(i)
    table := utils.SnakeCasedName(utils.StructName(i))
    r, err := db.Insert(table, m)

    if err == nil {
        insertId, err2 := r.LastInsertId()
        if err2 == nil && insertId > 0 {
            ps := reflect.ValueOf(i)
            if ps.Kind() == reflect.Ptr {
                // struct
                s := ps.Elem()
                if s.Kind() == reflect.Struct {
                    // exported field
                    f := s.FieldByName("Id")
                    if f.IsValid() {
                        // A Value can be changed only if it is 
                        // addressable and was not obtained by 
                        // the use of unexported struct fields.
                        if f.CanSet() {
                            // change value of N
                            if f.Kind() == reflect.Int {
                                if !f.OverflowInt(insertId) {
                                    f.SetInt(insertId)
                                }
                            }
                        }
                    }
                }
            }
        }
    }
    return r, err
}

func (db *DB) Update(table string, vals map[string]interface{}, where string, whereParams ...interface{}) (result sql.Result, err error) {
    if where == "" {
        panic("Can not update rows without where")
    }
    l := len(vals)
    if vals == nil || l < 1 {
        err = errors.New("no values to update")
        return
    }
    fields := make([]string, 0, l)
    params := make([]interface{}, 0, l)
    for k, v := range vals {
        switch v.(type) {
        case SQLLiteral:
            fields = append(fields, fmt.Sprintf("%s=%s", k, v))
        default:
            fields = append(fields, k+"=?")
            params = append(params, v)
        }
    }
    if whereParams != nil {
        params = append(params, whereParams...)
    }
    query := fmt.Sprintf("UPDATE %s SET %s WHERE %s;",
        table,
        strings.Join(fields, ", "),
        where)

    result, err = db.Exec(query, params...)
    return
}

func (db *DB) Delete(table string, where string, params ...interface{}) (result sql.Result, err error) {
    if where == "" {
        panic("Can not delete rows without where")
    }
    query := fmt.Sprintf("DELETE FROM %s WHERE %s;", table, where)

    result, err = db.Exec(query, params...)
    return
}

func (db *DB) rawSelectByStruct(structType reflect.Type, qi SqlQueryInfo) (rows *sql.Rows, fields []string, err error) {
    // nums of struct's fields
    lf := structType.NumField()
    // type's fields
    fields = make([]string, 0, lf)
    // sql select columns, it's Snake Cased
    columns := make([]string, 0, lf)

    // get fields in structType,
    // and convert to sql query column name
    for i := 0; i < lf; i++ {
        structField := structType.Field(i)
        fieldName := structField.Name
        fields = append(fields, fieldName)
        columns = append(columns, utils.SnakeCasedName(fieldName))
    }

    // tableName := utils.SnakeCasedName(utils.StructName(s))
    tableName := utils.SnakeCasedName(structType.Name())
    // TODO: check the fileds has specified ?
    qi.Fields = strings.Join(columns, ", ")
    // run query from db
    rows, err = db.Select(tableName, qi)
    return
}

// query by s and set the result value to s
// field mapping rule is: HelloWorld => hello_world
// mean that struct's field "HelloWorld" in database table's field is "hello_world"
// table name mapping use the same rule as field
func (db *DB) GetStruct(s interface{}, where string, params ...interface{}) error {

    structType := reflect.TypeOf(s)
    if structType.Kind() != reflect.Ptr {
        return errors.New(fmt.Sprintf("struct must be a pointer, but got %v", structType))
    }
    structType = structType.Elem()

    v := reflect.ValueOf(s)
    if v.IsNil() {
        return errors.New(fmt.Sprintf("struct can not be nil, but got %v", s))
    }

    qi := SqlQueryInfo{
        Limit:  1,
        Where:  where,
        Params: params,
    }
    rows, fields, err := db.rawSelectByStruct(structType, qi)
    if err != nil {
        return err
    }
    defer rows.Close()

    if rows.Next() {
        err = rawScanStruct(v, fields, rows)
        if err != nil {
            return err
        }
        if moreThanOneRow := rows.Next(); moreThanOneRow {
            return errors.New("more than one row found")
        }
    }
    return nil
}

// query by s and return a slice by type s
// field mapping rule is: HelloWorld => hello_world
// mean that struct's field "HelloWorld" in database table's field is "hello_world"
// table name mapping use the same rule as field
// @param s: just for reflect
// @return: notice that slice's item is pointer, like []*Blog
func (db *DB) GetStructs(s interface{}, qi SqlQueryInfo) ([]interface{}, error) {
    structType := reflect.TypeOf(s)
    if structType.Kind() == reflect.Ptr {
        structType = structType.Elem()
    }

    rows, fields, err := db.rawSelectByStruct(structType, qi)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    list := make([]interface{}, 0)

    for rows.Next() {
        v := reflect.New(structType)
        err = rawScanStruct(v, fields, rows)
        if err != nil {
            return nil, err
        }

        list = append(list, v.Interface())
    }

    return list, nil
}

// row scaner interface
// for sql.Row & sql.Rows
type rowScaner interface {
    Scan(dest ...interface{}) error
}

// scan the value by fields, and set to v
func rawScanStruct(v reflect.Value, fields []string, scanner rowScaner) (err error) {
    if v.IsNil() {
        e := fmt.Sprintf("struct can not be nil, but got %#v", v.Interface())
        return errors.New(e)
    }
    dest := make([]interface{}, len(fields))
    for v.Kind() == reflect.Ptr {
        v = v.Elem()
    }

    // Loop over column names and find field in s to bind to
    // based on column name. all returned columns must match
    // a field in the s struct
    for x, fieldName := range fields {
        f := v.FieldByName(fieldName)
        if f == zeroVal {
            e := fmt.Sprintf("Scanner: No field %s in type %s",
                fieldName, v.Type())
            return errors.New(e)
        } else {
            dest[x] = f.Addr().Interface()
        }
    }

    err = scanner.Scan(dest...)
    return
}

// mysql db
type MysqlDB struct {
    DB
}

// open mysql db, and return MysqlDB struct
func OpenMysql(driverName, dataSourceName string) (db *MysqlDB, err error) {
    var db2 *sql.DB
    db2, err = sql.Open(driverName, dataSourceName)
    if err != nil {
        return
    }
    db = &MysqlDB{}
    db.DB.DB = *db2
    return
}

// // from http://blog.csdn.net/hopingwhite/article/details/6107020
// var MySQLPool chan *mysql.MySQL
// func getMySQL() *mysql.MySQL {
//     if MySQLPool == nil {
//         MySQLPool = make(chan *mysql.MySQL, MAX_POOL_SIZE)
//     }
//     if len(MySQLPool) == 0 {
//         go func() {
//             for i := 0; i < MAX_POOL_SIZE/2; i++ {
//                 mysql := mysql.New()
//                 err := mysql.Connect("127.0.0.1", "root", "", "wgt", 3306)
//                 if err != nil {
//                     panic(err.String())
//                 }   
//                 putMySQL(mysql)
//             }   
//         }() 
//     }   
//     return <-MySQLPool
// }
// func putMySQL(conn *mysql.MySQL) {
//     if MySQLPool == nil {
//         MySQLPool = make(chan *mysql.MySQL, MAX_POOL_SIZE)
//     }   
//     if len(MySQLPool) == MAX_POOL_SIZE {
//         conn.Close()
//         return
//     }
//     MySQLPool <- conn
// }
