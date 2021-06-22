package main

import (
	"html/template"
	"os"
)

type User struct {
	Name  string
	Email string
	Age   int
}

func (u User) IsOld() bool {
	return u.Age > 40
}

func main() {
	user := User{Name: "hk", Email: "hkseo73@gmail.com", Age: 24}
	user2 := User{Name: "zzz", Email: "zzz", Age: 54}
	users := []User{user, user2}
	// {{}}는 빈공간을 뜻함 .은 Execute()로 넘어온 인스턴스를 뜻함
	tmpl, err := template.New("").ParseFiles("Templates/tmpl1.tmpl", "Templates/tmpl2.tmpl")
	if err != nil {
		panic(err)
	}
	tmpl.ExecuteTemplate(os.Stdout, "tmpl2.tmpl", users) // 탬플릿에 user정보를 채움
}
