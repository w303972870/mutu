package mtcore
import(
   "fmt"
   "time"
   "net/http"
)

type Server struct {}
type IServer interface{
    Start()
}

var HttpServer * Server

func( s * Server ) Start () {

    //新建一个路由管理器，
    mux := http.NewServeMux()
    //绑定路径和处理函数
    mux.HandleFunc("/test/usr/ddddd", htmlHandler)
    mux.Handle("/", http.HandlerFunc(indexHandler))

    server := & http.Server {
        Addr: MtTools.Str( 
        LibConfigParms.Configs.GetHttp( ConfigHttpIpKey ).(string) , ":" , 
        MtTools.Int2Str( LibConfigParms.Configs.GetHttp( ConfigHttpPortKey ).(int) ) ) ,
        Handler : mux ,
        ReadTimeout:    time.Duration( LibConfigParms.Configs.GetHttp( ConfigRTimeOutKey ).(int) ) * time.Second,
        WriteTimeout:   time.Duration( LibConfigParms.Configs.GetHttp( ConfigWTimeOutKey ).(int) ) * time.Second ,
        MaxHeaderBytes: LibConfigParms.Configs.GetHttp( ConfigHeaderBytesKey ).(int) ,
    }
    MutuLogs.Sys( MtTools.Str( 
        "启动HttpServer:" , LibConfigParms.Configs.GetHttp( ConfigHttpIpKey ).(string) , 
        ":" , MtTools.Int2Str( LibConfigParms.Configs.GetHttp( ConfigHttpPortKey ).(int) ) ) )

    err := server.ListenAndServe()
    if err != nil {
        MutuLogs.Error( err.Error() )
    }
}


func indexHandler( w http.ResponseWriter , r *http.Request) {    
    fmt.Fprintf(w, "hello world")
}

func htmlHandler(  w http.ResponseWriter , r *http.Request) {    
    w.Header().Set("Content-Type", "text/html")
    html := `<!doctype html>
             <META http-equiv="Content-Type" content="text/html" charset="utf-8">
             <html lang="zh-CN">
             <head>
                <title>Golang</title>
                <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=0;" />
             </head>
             <body>
                <div id="app">Welcome!</div>
             </body>
             </html>`
    fmt.Fprintf(w, html)
}


