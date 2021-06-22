package myapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var userMap map[int]*User
var lastId int

// User struct
type User struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	CreatedAt time.Time
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

// /users로 GET 요청 시 User 리스트 반환
func usersHandler(w http.ResponseWriter, r *http.Request) {
	if len(userMap) == 0 {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User")
		return
	}
	users := []*User{}
	for _, u := range userMap {
		users = append(users, u)
	}

	data, _ := json.Marshal(users)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(data))
}

// /users/{id}로 GET 요청 시 해당 유저 정보 반환
func getUserInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	user, ok := userMap[id]
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User ID:", id)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))
}

// /users로 유저 정보와 함께 POST 요청 시 유저 생성
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user) // 유저정보를 디코드하여 user객체에 입력
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	// Created User
	lastId++
	user.ID = lastId
	user.CreatedAt = time.Now()

	userMap[user.ID] = user // userMap에 저장

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	data, _ := json.Marshal(user) // json 형식으로 인코딩
	fmt.Fprint(w, string(data))
}

// /users/id로 DELETE 요청이 왔을 때 해당 유저 삭제
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)                 // url 변수 파싱 - gorilla mux 기능
	id, err := strconv.Atoi(vars["id"]) // 변수를 정수로 변환
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	_, ok := userMap[id] // 해당 아이디의 유저가 있다면
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User ID:", id)
		return
	}
	delete(userMap, id) // map 원소 삭제는 delete()
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Delete User ID:", id)
}

// /users로 업데이트 정보와 함께 PUT 요청이 왔을 때
func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	updateUser := new(User)
	err := json.NewDecoder(r.Body).Decode(updateUser) // 업데이트 정보를 updateUser에 저장
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	user, ok := userMap[updateUser.ID] // 해당 유저가 있다면
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User ID:", updateUser.ID)
		return
	}

	// 업데이트 정보가 있는 경우에는 그것으로 바꿔줌
	if updateUser.FirstName != "" {
		user.FirstName = updateUser.FirstName
	}
	if updateUser.LastName != "" {
		user.LastName = updateUser.LastName
	}
	if updateUser.Email != "" {
		user.Email = updateUser.Email
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))

}

// NewHandler make a new handler
func NewHandler() http.Handler {
	userMap = make(map[int]*User)
	lastId = 0
	mux := mux.NewRouter()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/users", usersHandler).Methods("GET") // 고릴라 먹스 덕에 같은 url에서도 모드에 따라 역할을 분리시킬 수 있다.
	mux.HandleFunc("/users", createUserHandler).Methods("POST")
	mux.HandleFunc("/users", updateUserHandler).Methods("PUT")
	mux.HandleFunc("/users/{id:[0-9]+}", getUserInfo).Methods("GET")
	mux.HandleFunc("/users/{id:[0-9]+}", deleteUserHandler).Methods("DELETE")
	return mux
}
