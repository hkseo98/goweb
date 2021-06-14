package myapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type User struct {
	FirstName string
	LastName  string
	Email     string
	CreatedAt time.Time
}

func indexHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprint(rw, "Hello Go!!!")
}

func barHandler(rw http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	fmt.Fprint(rw, "Hello "+name+"!")
}

type fooHandler struct {
	mu sync.Mutex // guards n
	n  int
}

func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user) // json 데이터를 request의 Body에서 받아서 원하는 형식으로 변환.
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request: ", err)
		return
	}
	user.CreatedAt = time.Now()
	data, _ := json.Marshal(user)                      // 데이터를 다시 json으로 변환
	w.Header().Add("Content-Type", "application/json") // 헤더에 콘텐츠 타입을 명시해줘야 json 형식인 걸 알아먹는다.
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(data))
}

// NewHttpHandler serve Mux
func NewHttpHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)

	mux.HandleFunc("/bar", barHandler)

	mux.Handle("/foo", new(fooHandler))

	http.ListenAndServe(":3000", mux)
	return mux
}
