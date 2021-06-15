package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadTest(t *testing.T) {
	assert := assert.New(t)
	path := "/Users/mac/Desktop/스크린샷 2021-06-15 오후 9.50.56.png"

	file, _ := os.Open(path)
	defer file.Close()

	os.RemoveAll("/uploads")

	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)                                  // multipart 파일을 쓰기 위한 writer
	w, err := writer.CreateFormFile("upload-file", filepath.Base(path)) // form에서 받는 형식
	assert.NoError(err)

	io.Copy(w, file) // 폼파일에 file을 복사
	writer.Close()

	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/uploads", buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	uploadsHandler(res, req)
	assert.Equal(http.StatusOK, res.Code)

	uploadFilePath := "./uploads/" + filepath.Base(path)
	_, error := os.Stat(uploadFilePath)
	assert.NoError(error)

	uploadFile, _ := os.Open(uploadFilePath)
	originFile, _ := os.Open(path)
	defer uploadFile.Close()
	defer originFile.Close()

	uploadData := []byte{}
	originData := []byte{}
	uploadFile.Read(uploadData)
	originFile.Read(originData)

	assert.Equal(uploadData, originData)
}
