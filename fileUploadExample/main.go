package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func uploadsHandler(w http.ResponseWriter, r *http.Request) {
	uploadfile, header, err := r.FormFile("upload-file") // input form으로 보낸 파일을 읽겠다
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	defer uploadfile.Close()

	dirname := "./uploads"
	os.MkdirAll(dirname, 0777)
	filepath := fmt.Sprintf("%s/%s", dirname, header.Filename)
	file, err := os.Create(filepath)
	defer file.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	io.Copy(file, uploadfile)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, filepath)
}

func main() {
	// public 폴더 안에 있는 파일들에 접근하는 서버 구동
	http.Handle("/", http.FileServer(http.Dir("public")))

	http.HandleFunc("/upload", uploadsHandler)

	http.ListenAndServe(":3000", nil)
}
