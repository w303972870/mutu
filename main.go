package main
import (
   "fmt"
   "framework/pkgs"
   //"framework/models"
   "reflect"
)

func main(){
    //fmt.Println(mutu.Conf().CArgs)
    //models.Qtest()
    //time.Sleep( time.Duration(5) * time.Second )
    type em struct{}

    fmt.Println(reflect.TypeOf(em{}).PkgPath())
    if mutu.Conf().CArgs.App == "http" {
        mutu.Server().Start()
    }
}

