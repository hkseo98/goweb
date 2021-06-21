package myapp

import (
	"encoding/json"
	"fmt"
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
	assert.Contains(string(data), "No User") // 앞에 있는 놈이 뒤에 있는 놈을 포함해야 OK

}

func TestGetUserInfo(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(NewHandler()) // 목업 테스트 서버 실행.
	defer ts.Close()

	res, err := http.Get(ts.URL + "/users/999")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	data, _ := ioutil.ReadAll(res.Body)
	assert.Contains(string(data), "No User ID:999")
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

	id := user.ID                                                 // 만들어진 유저의 아이디로 다시 GET 호출
	resp, errr := http.Get(ts.URL + "/users/" + strconv.Itoa(id)) // itoa = 정수를 문자열로
	assert.NoError(errr)
	assert.Equal(http.StatusOK, resp.StatusCode)
	user2 := new(User)
	e := json.NewDecoder(resp.Body).Decode(user2)
	assert.NoError(e)
	assert.Equal(user.ID, user2.ID) // 두 유저가 같다면 pass
	assert.Equal(user.FirstName, user2.FirstName)
}

func TestDeleteUser(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(NewHandler()) // 목업 테스트 서버 실행.
	defer ts.Close()

	req, _ := http.NewRequest("DELETE", ts.URL+"/users/1", nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	data, _ := ioutil.ReadAll(res.Body)
	// log.Println(string(data)) // log로 콘솔에 찍어볼 수 있음.
	assert.Contains(string(data), "No User ID:1")

	res, err = http.Post(ts.URL+"/users", "application/json", strings.NewReader(`{"FirstName": "hankyeol", "LastName":"seo", "Email":"hkseo73@gmail.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, res.StatusCode)

	user := new(User)
	err = json.NewDecoder(res.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	res, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	data, _ = ioutil.ReadAll(res.Body)
	// log.Println(string(data)) // log로 콘솔에 찍어볼 수 있음.
	assert.Contains(string(data), "Delete User ID:1")
}

func TestUpdateUser(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(NewHandler()) // 목업 테스트 서버 실행.
	defer ts.Close()

	req, _ := http.NewRequest("PUT", ts.URL+"/users",
		strings.NewReader(`{"ID":1, "FirstName":"updated", "LastName":"updated", "Email":"updated@naver.com"}`))
	res, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	data, _ := ioutil.ReadAll(res.Body)
	assert.Contains(string(data), "No User ID:1") // 없으면 업데이트 안 되게 하는 정책

	res, err = http.Post(ts.URL+"/users", "application/json", strings.NewReader(`{"FirstName": "hankyeol", "LastName":"seo", "Email":"hkseo73@gmail.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, res.StatusCode)

	user := new(User)
	err = json.NewDecoder(res.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	updateStr := fmt.Sprintf(`{"ID":%d, "FirstName":"updated"}`, user.ID)
	req, _ = http.NewRequest("PUT", ts.URL+"/users",
		strings.NewReader(updateStr))
	res, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	updatedUser := new(User)
	err = json.NewDecoder(res.Body).Decode(updatedUser)
	assert.NoError(err)

	assert.Equal(user.ID, updatedUser.ID)
	assert.Equal("updated", updatedUser.FirstName)
	assert.Equal(user.LastName, updatedUser.LastName)
	assert.Equal(user.Email, updatedUser.Email)
}

func TestUsers_WithData(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(NewHandler()) // 목업 테스트 서버 실행.
	defer ts.Close()

	res, err := http.Post(ts.URL+"/users", "application/json", strings.NewReader(`{"FirstName": "hankyeol", "LastName":"seo", "Email":"hkseo73@gmail.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, res.StatusCode)

	res, err = http.Post(ts.URL+"/users", "application/json", strings.NewReader(`{"FirstName": "hyun", "LastName":"ahn", "Email":"ahn@gmail.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, res.StatusCode)

	res, err = http.Get(ts.URL + "/users")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	users := []*User{}
	err = json.NewDecoder(res.Body).Decode(&users)
	assert.NoError(err)
	assert.Equal(2, len(users))
}
