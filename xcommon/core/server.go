package mtcore
import(
    "fmt"
    "time"
    "strings"
    "net/http"
    "io/ioutil"
    "bufio"
    "os"
    "regexp"
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
    actionfind(mux )
    //mrouter.AddRoute(mux)

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

func actionfind( mux * http.ServeMux ) {
    controllerPath := MtTools.Str( LibConfigParms.Configs.BasePath , "/" , LibConfigParms.CArgs.Package )
    if exist,_ := MtTools.PathExists( controllerPath ) ; exist == false {
        MutuLogs.Error( MtTools.Str( "项目控制器目录不存在，请检查：" , controllerPath ) )
    }
    txt := MtTools.Str( "package mrouter\nimport(\n\t\"fmt\"\n\t\"net/http\"\n\t\"framework/" , LibConfigParms.CArgs.Package,"\"\n)\nvar Router map[string]func( http.ResponseWriter , http.Request )\nfunc init() {\n\tRouter = make( map[string]func( http.ResponseWriter , http.Request ) , 0 )\n" )
    for _ , file := range MtTools.LsFiles( controllerPath ) {
        if file == "" {
            continue
        }
        for packageName , list := range tickFuncs(file , controllerPath) {
            for  controller , actionList := range list {
                for funcName,action := range actionList {
                    url := MtTools.Str( "/" , packageName , "/" , controller , "/" , action )
                    fmt.Println(url)
                    txt = MtTools.Str( txt , "\tRouter[\"", url , "\"]=", packageName , ".",funcName ,"\n" )
                }
            }
        }
    }
    txt = MtTools.Str( txt , "\n}")
    fmt.Println(txt)



    //fmt.Println(f)
    //for _, s := range f.Imports {
    //    fmt.Println(s.Path.Value)
    //}
    filePutContents( txt , MtTools.Str( LibConfigParms.Configs.BasePath , "/runtime/router.go" ) )
    //fmt.Println(txt)
}

/*
func actionfind1( mux * http.ServeMux ) {
    controllerPath := MtTools.Str( LibConfigParms.Configs.BasePath , "/" , LibConfigParms.CArgs.Package )
    if exist,_ := MtTools.PathExists( controllerPath ) ; exist == false {
        MutuLogs.Error( MtTools.Str( "项目控制器目录不存在，请检查：" , controllerPath ) )
    }
    pkg,err := importer.Default().Import("/home/wangdianchen/framework/api/")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(pkg)
    txt := MtTools.Str( "package mrouter\nimport(\n\"net/http\"\n\"framework/" , LibConfigParms.CArgs.Package,"\"\n)\nfunc AddRoute( mux * http.ServeMux ) {\n" )
    for _ , file := range MtTools.LsFiles( controllerPath ) {
        if file == "" {
            continue
        }
        for packageName , list := range tickFuncs(file , controllerPath) {
            for  controller , actionList := range list {
                for funcName,action := range actionList {
                    url := MtTools.Str( "/" , packageName , "/" , controller , "/" , action )
                    fmt.Println(url)
                    txt = MtTools.Str( txt , "mux.HandleFunc(\"", url , "\"",", " , packageName , ".",funcName ,")\n" )
                }
            }
        }
    }
    txt = MtTools.Str( txt , "}\n")

    fset := token.NewFileSet() 
    f, err := parser.ParseFile(fset, "", txt, parser.ParseComments )

    if err != nil {
        fmt.Println(err)
    }
    for _, i := range f.Decls {
        fn, ok := i.(*ast.FuncDecl)
        if !ok {
            continue
        }
        fmt.Println( fn.Name.Name )
    }

    ast.Inspect(f, func(n ast.Node) bool {
        call, ok := n.(*ast.CallExpr)
        if !ok {
            return true
        }

        fmt.Println(os.Stdout, fset, call.Fun)
        
        return false
    })
    for _, i := range f.Decls {
        fn, ok := i.(*ast.FuncDecl)
        if !ok {
            continue
        }
        //fn.Name.Obj.Decl.(ast.FuncDecl)(mux)
        fmt.Println("function:" )
    }

    //fmt.Println(f)
    //for _, s := range f.Imports {
    //    fmt.Println(s.Path.Value)
    //}
    //filePutContents( txt , MtTools.Str( LibConfigParms.Configs.BasePath , "/runtime/router.go" ) )
    //fmt.Println(txt)
}
*/
func makeMux( mux * http.ServeMux , url string , function string , packagename string ) {
    //mux.HandleFunc( url , reflect.Type.MethodByName( MtTools.Str( packagename , "." , function ) ) )
}

func tickFuncs( file string ,controllerPath string ) map[string]map[string]map[string]string {
        flysnowRegexp := regexp.MustCompile(`[A-Z][A-Za-z0-9]+Action`)
        var packageName string
        var funcList = make( map[string]map[string]map[string]string , 0 )

        controllerName := strings.TrimSuffix( file , `.go` )
        fileBuffer, err := os.Open( MtTools.Str( controllerPath ,"/", file )  )
        if err != nil {
            MutuLogs.Error( err.Error() )
        }
        scanner := bufio.NewScanner( fileBuffer )
        for scanner.Scan() {
            line := scanner.Text()  // or
            if len( line ) > 9 {
                if line[0:7] == "package" {
                    packageName = strings.TrimSpace( line[7:] )
                    funcList[packageName] = make( map[string]map[string]string , 0 )
                } else if line[0:4] == "func" {
                    funcName := flysnowRegexp.FindStringSubmatch( line )
                    if len( funcName ) > 0 {
                        if len( funcList[ packageName ][ controllerName ] ) == 0 {
                            funcList[ packageName ][ controllerName ] = make( map[string]string )
                        }
                        funcList[ packageName ][ controllerName ][ funcName[0] ] = getActionName( funcName[0] )
                    }
                }
            }
        }
        fileBuffer.Close()
        return funcList
}

func getActionName( funcName string ) string {
    return strings.ToLower( strings.TrimSuffix( funcName , `Action` ) )
}

func filePutContents( msg string , file string )  {
    if err := ioutil.WriteFile( file , []byte( msg ), 640 ); err != nil{
        MutuLogs.Error( err.Error() )
    }
}
