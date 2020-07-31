package mtcore
import(
    "fmt"
)

type iMtLogs interface {
    Waring( msg interface{} )
    Error( msg interface{} )
    Sys( msg interface{} )
    Info( msg interface{} )
}
type MtLogs struct{}

/*公用日志*/
var MutuLogs MtLogs

func( mt * MtLogs ) Waring ( message interface{} ) {
    fmt.Printf("%c[7;46;33m[警告]%s%c[0m\n", 0x1B, message.(string), 0x1B)
}

func( mt * MtLogs ) Error ( message interface{} ) {
    fmt.Printf("%c[5;41;32m[错误]%s%c[0m\n", 0x1B, message.(string), 0x1B)
    MtTools.Bye(1)
}

func( mt * MtLogs ) Sys ( message interface{} ) {
    fmt.Printf("%c[1;40;32m[系统]%s%c[0m\n", 0x1B, message.(string), 0x1B)
}

func( mt * MtLogs ) Info ( message interface{} ) {
    fmt.Println( "[信息]" , message.(string) )
}