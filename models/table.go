package models
import (
    "fmt"
    "framework/xcommon/core"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

/*Table定义*/
type Table struct {
    TbName string
    DbName string
}

type INametable interface {
    GetName() string
}

type Itable interface {
    Select( fields string , where *Condition ) *SelectResult
    SelectOne( fields string , where *Condition ) map[string]string
    Count( fields string , where *Condition ) int64
    Delete( where *Condition ) int64
}

func ( t * Table ) GetName() string {
    return t.TbName
}

/*条件定义*/
type Condition struct {
    Sql string
    Vals []interface{}
    LimitStr string
    OrderStr string
    GroupStr string
}

type ICondition interface {
    Set( sql string  , vals ...interface{} )
    Limit( start int , end int )
    GroupBy( by string )
    OrderBy( by string )
    GetSql() string
}

func( w * Condition )Set( sql string , vals ...interface{} ) {
    w.Sql = sql
    w.Vals = vals
}

func( w * Condition )Limit( start int , end int ) {
    if start > 0 && end > 0 {
        w.LimitStr = mtcore.MtTools.Str( " LIMIT " , mtcore.MtTools.Int2Str( start ) , " , " , mtcore.MtTools.Int2Str( end ) )
    } else {
        w.LimitStr = ""
    }
}

func( w * Condition )GroupBy( by string ) {
    if by != "" {
        w.GroupStr = mtcore.MtTools.Str( " GROUP BY " , by )
    } else {
        w.GroupStr = ""
    }
}

func( w * Condition )OrderBy( by string ) {
    if by != "" {
        w.OrderStr = mtcore.MtTools.Str( " ORDER BY " , by )
    } else {
        w.OrderStr = ""
    }
}

func( w * Condition )GetSql( ) string {
    sql := ""
    if w.Sql != "" {
        sql = mtcore.MtTools.Str( " WHERE " , w.Sql )
    }
    if w.GroupStr != "" {
        sql = mtcore.MtTools.Str( sql , w.GroupStr )
    }
    if w.OrderStr != "" {
        sql = mtcore.MtTools.Str( sql , w.OrderStr )
    }
    if w.LimitStr != "" {
        sql = mtcore.MtTools.Str( sql , w.LimitStr )
    }
    fmt.Println(sql)
    return sql
}

type SelectResult struct {
    Count int
    Rows []map[string]string
}

/*查询多条记录*/
func ( t * Table ) Select ( fields string , where *Condition ) *SelectResult {
    if fields == "" {
        fields = "*"
    }
    fmt.Println( fields ,t.DbName, where.Sql , where.Vals )
    sqlstr := mtcore.MtTools.Str( "select " , fields , " from " , t.GetName() , where.GetSql() )
    var result SelectResult
    mtcore.LibConfigParms.Mysqls[ t.DbName ].Connect( t.DbName )
    rows, err := mtcore.LibConfigParms.Mysqls[ t.DbName ].Db.Query( sqlstr , where.Vals... )

    if err != nil {
        mtcore.MutuLogs.Error( err.Error() )
        return & result
    }
    defer rows.Close()
    result.Rows = doRows( rows )
    result.Count = len( result.Rows )
    return & result
}

/*查询单条记录*/
func ( t * Table ) SelectOne ( fields string , where *Condition ) map[string]string {
    where.Limit( 0 , 1 )
    result := t.Select( fields , where )

    if result.Count < 1 {
        return make(map[string]string)
    } else {
        return result.Rows[0]
    }
}

/*查询记录数量*/
func ( t * Table ) Count ( where *Condition ,  fields string ) int64 {
    where.OrderBy("")
    where.GroupBy("")
    where.Limit(0,0)
    if fields == "" {
        fields = "1"
    }    
    result := t.SelectOne( mtcore.MtTools.Str( "count(" , fields , ") as total" ) , where )
    if len( result ) < 1 {
        return 0
    } else {
        return int64(mtcore.MtTools.Str2Int( result[ "total" ] ))
    }
}

/*删除记录*/
func ( t * Table ) Delete ( where *Condition ) int64 {
    where.OrderBy("")
    where.GroupBy("")
    where.Limit(0,0)
    sql := where.GetSql()
    if sql == "" || len(where.Vals) == 0 {
        mtcore.MutuLogs.Error( "为了数据安全，禁止没有where条件的DELETE" )
    }
    rows, err := mtcore.LibConfigParms.Mysqls[ t.DbName ].Db.Exec( mtcore.MtTools.Str( "DELETE FROM " , t.GetName() , sql ) , where.Vals... )
    if err != nil {
        mtcore.MutuLogs.Error( err.Error() )
    }
    affect , _ := rows.RowsAffected()
    return affect
}

func doRows( rows * sql.Rows )[]map[string]string {
    
    columns, err := rows.Columns()
    if err != nil {
        mtcore.MutuLogs.Error( err.Error() )
    }
    values := make( []sql.RawBytes , len( columns ) )
    scanArgs := make([]interface{}, len( values ) )

    var result []map[string]string

    for i := range values {
        scanArgs[i] = &values[i]
    }
    for rows.Next() {
        err = rows.Scan(scanArgs...)
        if err != nil {
            mtcore.MutuLogs.Error( err.Error() ) 
        }
        var value string
        one := make( map[string]string )

        for i, col := range values {
            if col == nil {
                value = ""
            } else {
                value = string(col)
            }
            one[ columns[i] ] = value
        }
        result = append( result ,  one )
    }
    if err = rows.Err(); err != nil {
        mtcore.MutuLogs.Error( err.Error() )
    }
    return result
}