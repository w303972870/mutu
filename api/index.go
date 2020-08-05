package api

import(
    "net/http"
    "fmt"
)

func IndexAction( w http.ResponseWriter , r *http.Request ) {
    fmt.Fprintf(w, "hello world index")
}

func HomeAction( w http.ResponseWriter , r *http.Request ) {
    fmt.Fprintf(w, "hello world home")
}


