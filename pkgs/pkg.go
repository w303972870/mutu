package mutu

import (
    "sync"
    "framework/xcommon/core"
)
var config * mtcore.ConfigParms
var logs * mtcore.MtLogs
var tools * mtcore.Tools
var once sync.Once

/*整理之后的系统变量*/
func Conf() * mtcore.ConfigParms {
    once.Do(func(){
        config = & mtcore.LibConfigParms
    })
    return config
}

/*整理之后的logs*/
func Logs() * mtcore.MtLogs {
    once.Do(func(){
        logs = & mtcore.MutuLogs
    })
    return logs
}

/*整理之后的tools*/
func Tools() * mtcore.Tools {
    once.Do(func(){
        tools = & mtcore.MtTools
    })
    return tools
}

func init(){

}
