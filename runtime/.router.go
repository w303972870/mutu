package mrouter
import(
"net/http"
"framework/api"
)
func AddRoute( mux * http.ServeMux ) {
mux.HandleFunc("/api/index/index", api.IndexAction)
mux.HandleFunc("/api/index/home", api.HomeAction)
}
