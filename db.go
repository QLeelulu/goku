package goku

import (
    "fmt"
    "database/sql"
    "strings"
    "errors"
    "reflect"
    "./utils"
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

type DB struct {
    sql.DB

    // Insert(table string, vals map[string]interface{}) (insertId int)
    // Update(table string, vals map[string]interface{}, where map[string]interface{}) (effectRow int)
    // Delete(table string, where map[string]interface{}) (effectRow int)
    // Select(table string, where map[string]interface{})
}

// // reture sql where query string
// //   where(map[string]interface{}{ "name":"lulu", "age": 23 }, " OR ")
// //   => "WHERE `name`=? OR `age`=?", []interface{}{ "lulu", 23 }
// func (db *DB) where(vars map[string]interface{}, grouping string) (sql string, params []interface{}) {
//     if vars == nil || len(vars) < 1 {
//         return
//     }
//     sqls = make([]string, 0, len(vars))
//     params = make([]interface{}, 0, len(vars))
//     for k, v := range vars {
//         sqls = append(sqls, fmt.Sprintf("`%s`=?", k))
//         params = append(params, v)
//     }
//     if grouping == "" {
//         grouping = " AND "
//     }
//     sql = strings.Join(sqls, grouping)
//     if sql {
//         sql = " WHERE " + sql
//     }
//     return
// }

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

    // get fields in s,
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
        v := reflect.ValueOf(s)

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
// @param s: just for reflect
// @return: notice that slice's item is pointer, like []*Blog
func (db *DB) GetStructs(s interface{}, qi SqlQueryInfo) ([]interface{}, error) {
    structType := reflect.TypeOf(s)
    if structType.Kind() == reflect.Ptr {
        structType = structType.Elem()
    }
    // // nums of struct's fields
    // lf := structType.NumField()
    // // type's fields
    // fields := make([]string, 0, lf)
    // // sql select columns, it's Snake Cased
    // columns := make([]string, 0, lf)

    // // get fields in s,
    // // and convert to sql query column name
    // for i := 0; i < lf; i++ {
    //     structField := structType.Field(i)
    //     fieldName := structField.Name
    //     fields = append(fields, fieldName)
    //     columns = append(columns, utils.SnakeCasedName(fieldName))
    // }

    // tableName := utils.SnakeCasedName(utils.StructName(s))
    // // TODO: check the fileds has specified ?
    // qi.Fields = strings.Join(columns, ", ")
    // // run query from db
    // rows, err := db.Select(tableName, qi)

    rows, fields, err := db.rawSelectByStruct(structType, qi)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    list := make([]interface{}, 0)
    //lf := len(fields)
    for rows.Next() {
        v := reflect.New(structType)
        // dest := make([]interface{}, lf)

        // // Loop over column names and find field in s to bind to
        // // based on column name. all returned columns must match
        // // a field in the s struct
        // for x, fieldName := range fields {
        //     f := v.Elem().FieldByName(fieldName)
        //     if f == zeroVal {
        //         e := fmt.Sprintf("db: No field %s in type %s (query: select %s)",
        //             fieldName, structType.Name(), utils.SnakeCasedName(structType.Name()))
        //         return nil, errors.New(e)
        //     } else {
        //         dest[x] = f.Addr().Interface()
        //     }
        // }

        // err = rows.Scan(dest...)

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
        e := fmt.Sprintf("value for %s can not be nil", v.Type())
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
