package main

import (
	"log"
	"net/http"
	"time"

	"github.io/hkseo98/goweb/WEBS/WebDecoratorHandler/decoHandler"
	"github.io/hkseo98/goweb/WEBS/WebDecoratorHandler/myapp"
)

// Decorator - 어떤 핸들러에 대한 요청이 있을 때 그 요청이 수행된 시간을 로그로 찍어주는 기능
func logger(w http.ResponseWriter, r *http.Request, h http.Handler) {
	start := time.Now()
	log.Println("[LOGGER1] Started")
	h.ServeHTTP(w, r) // mux가 호출될 때마다 logger가 실행됨
	log.Println("[LOGGER1] Completed time:", time.Since(start).Milliseconds())
}

func logger2(w http.ResponseWriter, r *http.Request, h http.Handler) {
	start := time.Now()
	log.Println("[LOGGER2] Started")
	h.ServeHTTP(w, r) // mux가 호출될 때마다 logger가 실행됨
	log.Println("[LOGGER2] Completed time:", time.Since(start).Milliseconds())
}

func NewHandler() http.Handler {
	mux := myapp.NewHandler()
	h := decoHandler.NewDecoHandler(mux, logger) // mux를 decorator함수인 logger로 감싼 것
	h = decoHandler.NewDecoHandler(h, logger2)   // h를 decorator함수인 logger2로 감싼 것
	// logger2가 가장 먼저 시작되고 가장 마지막에 끝남.
	return h
}

func main() {
	mux := NewHandler()
	http.ListenAndServe(":3000", mux)
}
