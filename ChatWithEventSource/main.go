package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/antage/eventsource"
	"github.com/gorilla/pat"
	"github.com/urfave/negroni"
)

func postMessageHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")
	name := r.FormValue("name")
	sendMessage(name, msg)
}

type Message struct {
	Name string
	Msg  string
}

var msgCh chan Message

func sendMessage(name, msg string) {
	// send message to every client
	msgCh <- Message{Name: name, Msg: msg}
}

func processMsgCh(es eventsource.EventSource) {
	for msg := range msgCh {
		data, _ := json.Marshal(msg)
		es.SendEventMessage(string(data), "", strconv.Itoa(time.Now().Nanosecond())) // es로 메세지 전송
	}
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("name")
	sendMessage("", "add user:"+username)
}

func leftUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("name")
	sendMessage("", "left user:"+username)
}

func main() {
	msgCh = make(chan Message)

	es := eventsource.New(nil, nil) // 이벤트소스 객체 만들기
	defer es.Close()

	go processMsgCh(es)

	mux := pat.New()
	mux.Post("/messages", postMessageHandler)
	mux.Handle("/stream", es)
	mux.Post("/users", addUserHandler)
	mux.Delete("/users", leftUserHandler)

	n := negroni.Classic()
	n.UseHandler(mux)
	http.ListenAndServe(":3000", n)
}
