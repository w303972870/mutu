package mutu
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
    Update( where *Condition ) int64
    Insert( insert  * Insert ) (int64,int64)
}

func ( t * Table ) GetName() string {
    return t.TbName
}

/*条件定义*/
type Condition struct {
    Sql string
    Vals []interface{}
    SetVals []interface{}
    SetStr string
    LimitStr string
    OrderStr string
    GroupStr string
}

type ICondition interface {
    Where( sql string  , vals ...interface{} )
    Limit( start int , end int )
    GroupBy( by string )
    OrderBy( by string )
    GetWhereSql() string
    GetSetSql() string
    GetSql() string
    Init()
}

func( w * Condition )Where( sql string , vals ...interface{} ) {
    w.Sql = sql
    w.Vals = vals
}

func( w * Condition )Set( sql string , vals ...interface{} ) {
    w.SetStr = sql
    w.SetVals = vals
}

func( w * Condition )Limit( start int , end int ) {
    if start >= 0 && end > 0 {
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
    return mtcore.MtTools.Str( w.GetSetSql() , w.GetWhereSql() )
}

func( w * Condition )GetWhereSql( ) string {
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
    return sql
}

func( w * Condition )GetSetSql() string {
    sql := ""
    if w.SetStr != "" {
        sql = mtcore.MtTools.Str( " SET " , w.SetStr )
    }
    return sql
}

func( w * Condition )Init( ) {
    w.OrderBy("")
    w.GroupBy("")
    w.Limit(0,0)
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
    sqlstr := mtcore.MtTools.Str( "select " , fields , " from " , t.GetName() , where.GetWhereSql() )

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
    where.Init()

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
    where.Init()

    sql := where.GetWhereSql()
    if sql == "" || len(where.Vals) == 0 {
        mtcore.MutuLogs.Error( "为了数据安全，禁止没有where条件的DELETE" )
    }

    rows, err := mtcore.LibConfigParms.Mysqls[ t.DbName ].Db.Exec( mtcore.MtTools.Str( "DELETE FROM " , t.GetName() , sql ) , where.Vals... )
    if err != nil {
        mtcore.MutuLogs.Error( err.Error() )
    }
    
    affect , err := rows.RowsAffected()

    return affect
}

/*更新记录*/
func ( t * Table ) Update ( where *Condition ) int64 {
    where.Init()
    sql := where.GetWhereSql()
    if sql == "" || len(where.Vals) == 0 {
        mtcore.MutuLogs.Error( "为了数据安全，禁止没有where条件的Update" )
    }

    sql = where.GetSql()
    vals := mtcore.MtTools.Merger( where.SetVals , where.Vals )
    rows, err := mtcore.LibConfigParms.Mysqls[ t.DbName ].Db.Exec( mtcore.MtTools.Str( "UPDATE " , t.GetName() , sql ) , vals... )

    if err != nil {
        mtcore.MutuLogs.Error( err.Error() )
    }
    affect , err := rows.RowsAffected()
    if err != nil {
        mtcore.MutuLogs.Error( err.Error() )
    }
    return affect
}

/*插入记录，返回值总数，最后插入id*/
func( t * Table )Insert( insert * Insert ) (int64,int64){
    checkInsert( insert )
    sql,vals := insert.GetSql()
    result, err := mtcore.LibConfigParms.Mysqls[ t.DbName ].Db.Exec( mtcore.MtTools.Str( "INSERT INTO " , t.GetName() , *sql ), *vals... )
    if err != nil {
        mtcore.MutuLogs.Error( err.Error() )
    }
    count , _ := result.RowsAffected()
    id , _ := result.LastInsertId()

    fmt.Println(id)
    return int64( count ),int64( id )
}


/*插入定义*/
type Insert struct {
    fields []string
    vals [][]interface{}
}

type IInsert interface {
    Fields( fields ...string )
    GetFields() []string
    Values( fields ...interface{} )
    GetValues()[][]interface{}
    GetSql() ( *string , *[]interface{} )
    Init()
}

/*设置插入字段列表*/
func ( t * Insert ) Fields ( fields ...string ) {
    t.fields = fields
}

func ( t * Insert ) Values ( vals ...interface{} ) {
    t.vals = append( t.vals , vals )
}

func ( t * Insert ) GetFields()[]string {
    return t.fields
}

func ( t * Insert ) GetValues()[][]interface{} {
    return t.vals
}

func ( t * Insert ) Init() {
    t.fields = make( []string , 0 )
    t.vals = make( [][]interface{} , 0 )
}

/*组织insert语句,返回语句和值*/
func ( t * Insert ) GetSql() ( *string , *[]interface{} ) {
    inserSql := " ( "
    valueSql := " ( "

    total := len( t.fields )
    for i,v := range t.fields {
        if total - 1 == i {
            inserSql = mtcore.MtTools.Str( inserSql , "`" , v , "`" , " ) VALUES " )
            valueSql = mtcore.MtTools.Str( valueSql , "?" , " ) " )
        } else {
            inserSql = mtcore.MtTools.Str( inserSql , "`" , v , "`" , " , " )
            valueSql = mtcore.MtTools.Str( valueSql , "?" , " , " )
        }
    }
    total = len( t.vals )
    var parms []interface{}

    for i,v := range t.vals {
        if total - 1 == i {
            inserSql = mtcore.MtTools.Str( inserSql , valueSql )
        } else {
            inserSql = mtcore.MtTools.Str( inserSql , valueSql , "," )
        }
        parms = append( parms , v... )
    }
    return &inserSql , &parms
}

/*检查insert合法性*/
func checkInsert( insert * Insert ) {
    if len(insert.GetValues()) == 0 || len(insert.GetFields()) == 0 {
        mtcore.MutuLogs.Error( "设置插入数据与字段不完整" )
    }

    for _,v := range insert.GetValues() {
        if len(v) != len(insert.GetFields()) {
            mtcore.MutuLogs.Error( "插入数据与字段长度不一致" )
            break
        }
    }
}

/*整理行数据*/
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