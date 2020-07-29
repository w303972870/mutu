package mtcore
import(
    "fmt"
    "sync"
)

type iMtLogs interface {
    Start()
    Waring( msg interface{} )
    Error( msg interface{} )
    Sys( msg interface{} )
    Info( msg interface{} )
}
type MtLogs struct{}

/*公用日志*/
var MutuLogs MtLogs

func( mt MtLogs ) Waring ( msg interface{} ) {
    LibConfigParms.LibLogChan.ChanWaring <- msg
}

func( mt MtLogs ) Error ( msg interface{} ) {
    LibConfigParms.LibLogChan.ChanError <- msg
}

func( mt MtLogs ) Sys ( msg interface{} ) {
    LibConfigParms.LibLogChan.ChanSys <- msg
}

func( mt MtLogs ) Info ( msg interface{} ) {
    LibConfigParms.LibLogChan.ChanInfo <- msg
}

func( mt MtLogs ) Start () {
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()
        for true {
            select {
                case msg , ok := <- LibConfigParms.LibLogChan.ChanInfo :
                    if ok {
                        showInfo( msg )
                    }
                case msg , ok := <- LibConfigParms.LibLogChan.ChanSys :
                    if ok {
                        showSys( msg )
                    }
                case msg , ok := <- LibConfigParms.LibLogChan.ChanError :
                    if ok {
                        showError( msg )
                    }
                case msg , ok := <- LibConfigParms.LibLogChan.ChanWaring :
                    if ok {
                        showWaring( msg )
                    }
            }
            continue
        }
    }()
    //wg.Wait()    
}

func showWaring ( message interface{} ){
    fmt.Printf("%c[7;46;33m[警告]%s%c[0m\n", 0x1B, message.(string), 0x1B)
}

func showError ( message interface{} ){
    fmt.Printf("%c[5;41;32m[错误]%s%c[0m\n", 0x1B, message.(string), 0x1B)
    MtTools.Bye(1)
}

func showSys ( message interface{} ){
    fmt.Printf("%c[1;40;32m[系统]%s%c[0m\n", 0x1B, message.(string), 0x1B)
}

func showInfo ( message interface{} ){
    fmt.Println( "[信息]" , message.(string) )
}

func init(){
    MutuLogs.Start()
}

