package mtcore

/*
   初始化参数
*/
import (
    "os"
    "fmt"
    "flag"
    "os/signal"
    "syscall"
    "runtime"
)

func init(){
    runtime.GOMAXPROCS( runtime.NumCPU() * 5 )

    flag.StringVar( &LibConfigParms.CArgs.App , "app" , "http" , "启动对象[http|cmd]" )
    flag.IntVar(&LibConfigParms.CArgs.Port, "P", 80, "当是http时的监听端口号")
    flag.StringVar(&LibConfigParms.CArgs.ConfigFile, "config", "./config/", "config.yaml配置文件的所在【目录】")

    flag.Usage = usage
    flag.Parse()
    if ! MtTools.Exist( LibConfigParms.CArgs.ConfigFile ) {
        MutuLogs.Error( "config error" )
        usage()
    }
    readConfig( "config" , LibConfigParms.CArgs.ConfigFile )

    setupSignalHandler()
}

//返回命令行参数
func Args() interface{} {
    fmt.Println(LibConfigParms.Configs.GetDatabase("face_detect","port"),LibConfigParms.Configs.GetCustom("logs"))
    return LibConfigParms.CArgs
}

//定义Usage样式
func usage() {
    fmt.Fprintf(os.Stderr, `framework version: mutu/1.0.0
Usage: mutu [-hvVtTq] [-s signal] [-c filename] [-p prefix] [-g directives]

Options:
`)
    flag.PrintDefaults()
    MtTools.Bye(0)
}

//定义control+c退出提示
func setupSignalHandler() {
    shutdownSignals := []os.Signal{os.Interrupt, syscall.SIGTERM}
	c := make(chan os.Signal, 2) 
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
        fmt.Println("再按一次Control+C退出...")
		<-c
		fmt.Println("\nControl+C退出....")
		MtTools.Bye(0)
	}()
}

