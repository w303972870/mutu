package main
import (
   //"fmt"
   "framework/pkgs"
   //"framework/models"

)


func main(){
    if mutu.Conf().CArgs.App == "http" {
        mutu.Server().Start()
    }
}
