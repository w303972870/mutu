package mrouter
import(
	"fmt"
	"net/http"
	"framework/api"
)
var Router map[string]func( http.ResponseWriter , http.Request )
func init() {
	Router = make( map[string]func( http.ResponseWriter , http.Request ) , 0 )
	Router["/api/index/index"]=api.IndexAction
	Router["/api/index/home"]=api.HomeAction

}