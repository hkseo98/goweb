package myapp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(NewHandler()) // 목업 테스트 서버 실행.
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	data, _ := ioutil.ReadAll(res.Body)
	assert.Equal("Hello World", string(data))
}

func TestUsers(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(NewHandler()) // 목업 테스트 서버 실행.
	defer ts.Close()

	res, err := http.Get(ts.URL + "/users")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	data, _ := ioutil.ReadAll(res.Body)
	assert.Contains(string(data), "Get UserInfo") // 앞에 있는 놈이 뒤에 있는 놈을 포함해야 OK
}

func TestGetUserInfo(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(NewHandler()) // 목업 테스트 서버 실행.
	defer ts.Close()

	res, err := http.Get(ts.URL + "/users/999")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	data, _ := ioutil.ReadAll(res.Body)
	assert.Contains(string(data), "User Id:999")
}

func TestCreateUserInfo(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(NewHandler()) // 목업 테스트 서버 실행.
	defer ts.Close()

	res, err := http.Post(ts.URL+"/users", "application/json", strings.NewReader(`{"FirstName": "hankyeol", "LastName":"seo", "Email":"hkseo73@gmail.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, res.StatusCode)

	user := new(User)
	er := json.NewDecoder(res.Body).Decode(user)
	assert.NoError(er)
	assert.NotEqual(0, user.ID)

	id := user.ID
	resp, errr := http.Get(ts.URL + "/users/" + strconv.Itoa(id)) // itoa = 정수를 문자열로
	assert.NoError(errr)
	assert.Equal(http.StatusOK, resp.StatusCode)
	user2 := new(User)
	e := json.NewDecoder(resp.Body).Decode(user2)
	assert.NoError(e)
	assert.Equal(user.ID, user2.ID)
	assert.Equal(user.FirstName, user2.FirstName)
}
