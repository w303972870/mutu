package mtcore
import(
    "os"
    "bytes"
    "reflect"
    "strconv"
)

type Tools struct{}

type itools interface {
    Exist( path string ) bool
    GetType( parm interface{} ) interface{}
    Str( parms ...string ) string
    Str2Int( str string ) int
    Int2Str( i int ) string
    Bye( code int )
}

var MtTools Tools 


/*判断文件路径是否存在
func Exist( path string ) bool {
    if LibConfigParms.SysType == SysLinux {
        err := syscall.Access( path, syscall.F_OK )
        return ! os.IsNotExist(err)        
    } else {
        _, err := os.Lstat(path)
        return ! os.IsNotExist(err)
    }
}
*/

func( t * Tools )Exist( path string ) bool {
    _, err := os.Lstat(path)
    return ! os.IsNotExist(err)
}

/*获取变量类型*/
func( t * Tools )GetType( parm interface{} ) interface{} {
    return reflect.TypeOf( parm )
}

/*拼接字符串*/
func( t * Tools )Str( parms ...string ) string {
    var buffer bytes.Buffer

    for _ , parm := range parms {
        buffer.WriteString( parm )
    }
    return buffer.String()
}

/*string转int*/
func( t * Tools )Str2Int( str string ) int {
    sint, err := strconv.Atoi( str )
    if err != nil {
        return -1
    } else {
        return sint
    }
}

/*int转string*/
func( t * Tools )Int2Str( i int ) string {
    return strconv.Itoa( i )
}

/*退出程序*/
func( t * Tools )Bye( code int ){
    os.Exit( code )
}