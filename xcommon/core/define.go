package mtcore

import(
    "runtime"
)

const (
    SysLinux = "Linux"
    SysWindows = "Windows"

    ConfigBasePathKey = "basepath"
    ConfigHttpKey = "server"
    ConfigHttpIpKey = "ip"
    ConfigHttpPortKey = "port"
    ConfigRTimeOutKey = "read_time_out"
    ConfigWTimeOutKey = "write_time_out"
    ConfigHeaderBytesKey = "max_header_bytes"
    ConfigDbKey = "database"
    ConfigDbNameKey = "db"
    ConfigDbHostKey = "host"
    ConfigDbPortKey = "port"
    ConfigDbUserKey = "user"
    ConfigDbPwdKey = "password"
    ConfigCustomKey = "custom"
)

/*整体系统配置变量信息*/
type ConfigParms struct {

    /*系统平台类型*/
    SysType string

    /*命令行参数*/
    CArgs struct {
        App string 
        ConfigFile string 
        Port int 
    }

    /*用户配置文件内容*/
    Configs *configs

    /*用户数据库对象*/
    Mysqls map[string] *Mysql
}

var LibConfigParms ConfigParms

/*定义配置内容接口*/
type iConfig interface {
    SetBasePath( path string )
    GetBasePath() string
    SetCustom( key string , value interface{} )
    GetCustom( key string ) interface{}
    SetHttp( key string , value interface{} )
    GetHttp( key string ) interface{}
    SetDatabase( database string , key string , value interface{} )
    GetDatabase( database string , key string ) interface{}
}

/*设置Http配置*/
func ( config * configs ) SetHttp( key string , value interface{} ) {
    if len( config.Http ) < 1 {
        config.Http = make( map[string]interface{} )
    }
    
    switch key {
        case ConfigHttpIpKey :
            config.Http[ key ] = value
        case ConfigHttpPortKey :
            config.Http[ key ] = value
        case ConfigRTimeOutKey :
            config.Http[ key ] = value
        case ConfigWTimeOutKey :
            config.Http[ key ] = value
        case ConfigHeaderBytesKey :
            if value != nil && value != "" && MtTools.GetType( value ) != "string" && value.(int) > 1 {
                config.Http[ key ] = value.(int) * 1024
            }
        default :
            //MutuLogs.Error( MtTools.Str( "配置文件错误： " , ConfigHttpKey , "," , key ) )
    }
}

/*获取Http配置*/
func ( config * configs ) GetHttp( key string ) interface{} {
    switch key {
        case ConfigHttpIpKey :
            if len( config.Http ) == 0 || config.Http[ key ] == "" || config.Http[ key ] == nil {
                //MutuLogs.Sys( "ServerIP错误，启用默认值\"0.0.0.0\"" )
                config.SetHttp( key , "0.0.0.0" )
            }
        case ConfigHttpPortKey :
            if len( config.Http ) == 0 || 
                config.Http[ key ] == "" || 
                MtTools.GetType( config.Http[ key ] ) == "string" || 
                config.Http[ key ].(int) < 1 {
                //MutuLogs.Sys( "ServerPort错误，启用默认值\"666\"" )
                config.SetHttp( key , 666 )
            }
        case ConfigRTimeOutKey :
            if len( config.Http ) == 0 || 
                config.Http[ key ] == "" || 
                MtTools.GetType( config.Http[ key ] ) == "string" || 
                config.Http[ key ].(int) < 1 {
                config.SetHttp( key , 10 )
            }
        case ConfigWTimeOutKey :
            if len( config.Http ) == 0 || 
                config.Http[ key ] == "" || 
                MtTools.GetType( config.Http[ key ] ) == "string" || 
                config.Http[ key ].(int) < 1 {
                config.SetHttp( key , 10 )
            }
        case ConfigHeaderBytesKey :
            if len( config.Http ) == 0 || 
                config.Http[ key ] == "" || 
                MtTools.GetType( config.Http[ key ] ) == "string" || 
                config.Http[ key ].(int) < 1 {
                config.SetHttp( key , 1024 )
            }
        default :
            MutuLogs.Error( MtTools.Str( "不存在此键： " , key ) )
    }
    return config.Http[ key ]
}

/*获取自定义配置*/
func ( config configs ) GetCustom( key string ) interface{} {
    return config.Custom[ key ]
}

/*设置自定义配置*/
func ( config * configs ) SetCustom( key string , value interface{} ) {
    if len( config.Custom ) < 1 {
        config.Custom = make( map[string]interface{} )
    }
    config.Custom[ key ] = value
}

/*设置程序路径*/
func ( config * configs ) SetBasePath( path string ) {
    config.BasePath = path
}

/*获取程序路径*/
func ( config configs ) GetBasePath() string {
    return config.BasePath
}

/*设置数据库配置*/
func ( config * configs ) SetDatabase( db string , key string , value interface{} ) {
    if _, ok := config.Database[ db ] ; !ok {
        config.Database[ db ] = new( cdb )
    }
    switch key {
        case ConfigDbHostKey :            
            config.Database[ db ].db = db
            config.Database[ db ].host = value.(string)
        case ConfigDbPortKey :
            config.Database[ db ].port = value.(int)
        case ConfigDbUserKey :
            config.Database[ db ].user = value.(string)
        case ConfigDbPwdKey :
            config.Database[ db ].password = value.(string)
        default :
            MutuLogs.Error( MtTools.Str( "配置文件错误： " , key ) )
    }
}

/*获取数据库配置*/
func ( config configs ) GetDatabase( db string , key string ) interface{} {
    switch key {
        case ConfigDbNameKey : 
            return config.Database[ db ].db
        case ConfigDbHostKey :            
            return config.Database[ db ].host 
        case ConfigDbPortKey :
            return config.Database[ db ].port 
        case ConfigDbUserKey :
            return config.Database[ db ].user 
        case ConfigDbPwdKey :
            return config.Database[ db ].password 
        default :
            MutuLogs.Error( MtTools.Str( "Database:" , db , " Not Found" ) )
            
        return nil
    }
}

/*配置文件*/
type configs struct {
    BasePath string
    Database map[string] * cdb
    Http map[string]interface{}
    Custom map[string]interface{}
}

/*配置文件-数据库*/
type cdb struct {
    db string
    host string
    port int
    user string
    password string
}

/*初始化*/
func init(){
    LibConfigParms.Configs = new( configs )
    LibConfigParms.Configs.Database = make( map[string] *cdb )
    LibConfigParms.Mysqls = make( map[string] *Mysql )
    LibConfigParms.Configs.Custom = make( map[string]interface{} )

    if runtime.GOOS == "linux" {
        LibConfigParms.SysType = SysLinux
    } else if runtime.GOOS == "windows" {
        LibConfigParms.SysType = SysWindows
    } else {
        LibConfigParms.SysType = "Unknown"
    }

    MutuLogs.Sys( MtTools.Str( "当前运行平台： " , LibConfigParms.SysType ) )
}

type Itable interface {
    GetName() string
}
