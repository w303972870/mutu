package mtcore
import(
    "reflect"
	"github.com/spf13/viper"
	"github.com/fsnotify/fsnotify"
)

func readConfig( configname string , configfile string ) {
	viper.SetConfigName( configname )
	viper.AddConfigPath( configfile )
	err := viper.ReadInConfig()
	if err != nil {
		MutuLogs.Error("读取不到配置文件")
	}
	viper.WatchConfig()
    viper.OnConfigChange(func(e fsnotify.Event) {
		MutuLogs.Sys( MtTools.Str( "配置发生变更：" , e.Name ) )
        initConfigs()
    })
    initConfigs()

}

func initConfigs(){
    LibConfigParms.Configs.SetBasePath( viper.Get( ConfigBasePathKey ).(string) )
    
    database := viper.GetStringMap( ConfigDbKey )
    if len(database) > 0 {
        readDbconfig( database )
    } else {
        MutuLogs.Waring( "未发现数据库配置" )
    }

    custom := viper.GetStringMap( ConfigCustomKey )
    if len(custom) > 0 {
        readCustomConfig( custom )
    }

    http := viper.GetStringMap( ConfigHttpKey )
    if len(http) > 0 {
        readHttpConfig( http )
    }
}

/*处理配置文件-http部分*/
func readHttpConfig( http map[string]interface{} ) {
    for k,v := range http {
        LibConfigParms.Configs.SetHttp( k , v )
    }
}

/*处理配置文件-自定义配置部分*/
func readCustomConfig( custom map[string]interface{} ) {
    for key,value := range custom {
        LibConfigParms.Configs.SetCustom( key , value )
    }
}

/*处理配置文件-数据库部分*/
func readDbconfig( database map[string]interface{} ) {
    for dbName,list := range database {
        for key,value := range list.(map[string]interface{}) {
            LibConfigParms.Configs.SetDatabase( dbName , key , value )
        }
        LibConfigParms.Mysqls[ dbName ] = new ( Mysql )
        LibConfigParms.Mysqls[ dbName ].Connect( dbName )
    }
}

func StructToMapDemo(obj interface{}) map[string]interface{}{
    obj1 := reflect.TypeOf(obj)
    obj2 := reflect.ValueOf(obj)
 
    var data = make(map[string]interface{})
    for i := 0; i < obj1.NumField(); i++ {
        data[obj1.Field(i).Name] = obj2.Field(i).Interface()
    }
    return data
}

// 遍历struct并且自动进行赋值
func structByReflect(data map[string]interface{}, inStructPtr interface{}) {
    rType := reflect.TypeOf(inStructPtr)
    rVal := reflect.ValueOf(inStructPtr)
    if rType.Kind() == reflect.Ptr {
        // 传入的inStructPtr是指针，需要.Elem()取得指针指向的value
        rType = rType.Elem()
        rVal = rVal.Elem()
    } else {
        panic("inStructPtr must be ptr to struct")
    }
    // 遍历结构体
    for i := 0; i < rType.NumField(); i++ {
        t := rType.Field(i)
        f := rVal.Field(i)
        if v, ok := data[t.Name]; ok {
            f.Set(reflect.ValueOf(v))
        } else {
            panic(t.Name + " not found")
        }
    }
}