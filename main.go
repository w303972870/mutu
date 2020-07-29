package main
import (
   "fmt"
   "time"
   "framework/pkgs"
   "framework/models"
)

func main(){
    fmt.Println(mutu.Conf().CArgs)
    models.Qtest()
    time.Sleep( time.Duration(5) * time.Second )
}