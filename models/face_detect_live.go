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
    where.Where( "live_id > ? and live_id < ?" , 38444 , 38546 )
    where.Limit( 0 , 10 )
    where.OrderBy( "live_id asc" )
    where.GroupBy( "live_id" )

    result := LiveTable.Select( "live_id,MAC,face_image" , &where )
    fmt.Println( result.Count )
    for r,l := range result.Rows {
        fmt.Println(r,l)
    }

    /*    

    resultone := LiveTable.SelectOne( "live_id,MAC,face_image" , &where )
    fmt.Println( resultone )

    //resultcount := LiveTable.Count( &where , "" )
    resultcount := LiveTable.Count( &where , "distinct MAC" )
    fmt.Println( resultcount )

    where.Where( "live_id > ? and live_id < ?" , 38444 , 38446 )
    delcount := LiveTable.Delete( &where )
    fmt.Println( delcount )
 
    where.Set( "MAC = ? , face_coordinate = ?" , "00:11:04:02:33:d6" , "{\"pingmian\":[0,0,0]}" )
    updatecount := LiveTable.Update( &where )
    fmt.Println( updatecount )

    where.Where( "live_id > ? " , 128893 )
    delcount := LiveTable.Delete( &where )
    fmt.Println( delcount )
 
    var insert Insert
    insert.Fields( "MAC" , "face_id" , "live_type" )
    insert.Values( "00:11:04:02:33:d6","1584058701053069824",1 )

    count,insertid := LiveTable.Insert( &insert )
    fmt.Println(count, insertid )
   */

}

