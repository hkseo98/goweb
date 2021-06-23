package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/pat"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var rd *render.Render

type User struct {
	Name      string
	Email     string
	CreatedAt time.Time
}

func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	user := User{Name: "hk", Email: "hkseo73@gmail.com"}
	rd.JSON(w, http.StatusOK, user)
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		rd.Text(w, http.StatusBadRequest, err.Error()) // err.Error() - 에러 메세지 반환
		return
	}
	user.CreatedAt = time.Now()

	rd.JSON(w, http.StatusOK, user)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// tmpl, err := template.New("").ParseFiles("templates/hello.tmpl")
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	fmt.Fprint(w, err)
	// 	return
	// }
	// tmpl.ExecuteTemplate(w, "hello.tmpl", "hkseo")

	user := User{Name: "hk", Email: "hkseo73@gmail.com"}
	rd.HTML(w, http.StatusOK, "body", user) // 한 줄로 끝남..
	// 템플릿은 반드시 templates라는 이름의 폴더 안에 있어야 함.
	// 다른 폴더를 참조하고 싶다면 rd 초기화 시 render.Option{}에서 설정해야 함
	// 템플릿의 확장명은 반드시 .tmpl이어야 한다.
	// 다른 확장명도 읽고 싶으면 rd 초기화 시 render.Option{}에서 설정해야 함.
	// 템플릿의 확장명은 함수에 적지 말아야 한다.

}

func main() {
	rd = render.New(render.Options{
		Directory:  "templ",                    // 폴더 변경!
		Extensions: []string{".html", ".tmpl"}, // 확장자 옵션에 html 추가
		Layout:     "hello",
	})
	mux := pat.New() // gorilla pat

	mux.Get("/users", getUserInfoHandler)
	mux.Post("/users", addUserHandler)
	mux.Get("/hello", helloHandler)

	// mux.Handle("/", http.FileServer(http.Dir("public")))
	n := negroni.Classic() // negroni는 기본적인 미들웨어를 제공해줌
	n.UseHandler(mux)      // mux를 감싸는 데코레이터, 로그, 파일서버 기능 등등

	http.ListenAndServe(":3000", n)
}
