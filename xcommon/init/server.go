package minit
import(
    "fmt"
    "io/ioutil"
)

func Test(){
    fmt.Println("test")
}

func init(){
    makeFile()
}

func makeFile(){
    str := "package mrouter\nimport(\n\"net/http\"\n\"framework/api\"\n)\nfunc AddRoute( mux * http.ServeMux ) {\nmux.HandleFunc(\"/api/index/index\", api.IndexAction)\nmux.HandleFunc(\"/api/index/home\", api.HomeAction)\n}"

    filePutContents( str , "/home/wangdianchen/framework/runtime/router.go" )

}

func filePutContents( msg string , file string )  {
    if err := ioutil.WriteFile( file , []byte( msg ), 640 ); err != nil{
        fmt.Println(err.Error())
    }
}