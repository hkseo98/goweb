package main

import (
	"net/http"

	"github.io/hkseo98/goweb/myapp"
)

func main() {
	http.ListenAndServe(":3000", myapp.NewHttpHandler())
}
