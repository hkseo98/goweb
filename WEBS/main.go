package main

import (
	"net/http"

	"github.io/hkseo98/goweb/WEBS/myapp.go"
)

func main() {
	http.ListenAndServe(":3000", myapp.NewHandler())
}
