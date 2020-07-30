package models
import (
    "fmt"
)

var LiveTable Table
func init(){
    //LiveTable.TbName = "live_data"
    LiveTable.TbName = "test"
    LiveTable.DbName = "face_detect_live"
}


func Qtest(){
    var where Condition
    where.Set( "live_id > ? and live_id < ?" , 38444 , 38544 )
    where.Limit( 0 , 10 )
    where.OrderBy( "live_id asc" )
    where.GroupBy( "live_id" )

    /*    
    result := LiveTable.Select( "live_id,MAC,face_image" , &where )
    fmt.Println( result.Count )
    for r,l := range result.Rows {
        fmt.Println(r,l)
    }

    resultone := LiveTable.SelectOne( "live_id,MAC,face_image" , &where )
    fmt.Println( resultone )

    //resultcount := LiveTable.Count( &where , "" )
    resultcount := LiveTable.Count( &where , "distinct MAC" )
    fmt.Println( resultcount )
    */
    where.Set( "live_id > ? and live_id < ?" , 38444 , 38446 )
    delcount := LiveTable.Delete( &where )
    fmt.Println( delcount )


}

