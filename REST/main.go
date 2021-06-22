package main

import (
	"net/http"

	"github.io/hkseo98/goweb/REST/myapp"
)

func main() {
	http.ListenAndServe(":3000", myapp.NewHandler())
}
