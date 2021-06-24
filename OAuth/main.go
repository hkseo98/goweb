package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/pat"
	"github.com/urfave/negroni"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// 구글 API 사이트에서 프로젝트를 생성하고 URL및, 리다이렉트 URL을 등록하고 clientid, secretkey를 생성해야 한다.

var googleOauthConfig = oauth2.Config{
	RedirectURL:  "http://localhost:3000/auth/google/callback", // 어디로 콜백을 알려줄 것인지 등록
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_SECRET_KEY"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"}, // 접근 범위 설정
	Endpoint:     google.Endpoint,
}

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	// AuthCodeURL는 state는 매개변수로 가짐 - State is a token to protect the user from CSRF attacks. You must always provide a non-empty string and validate that it matches the the state query parameter on your redirect callback.
	state := generateStateOauthCookie(w)
	// AuthCodeURL returns a URL to OAuth 2.0 provider's consent page that asks for permissions for the required scopes explicitly.
	url := googleOauthConfig.AuthCodeURL(state)
	// 구글 로그인 사이트로 리다이렉트 시켜줌
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	// 로그인이 끝나고 콜백이 불렸을 때 우리가 저장했던 쿠키와 구글에서 알려준 state 값을 비교해서 같으면 정상적 로그인임을 알 수 있음
}

// 브라우저에 쿠키에 템포러리 키를 심고 리다이렉트로 콜백이 왔을 때 그 쿠키를 다시 비교하는 방식
func generateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(1 * 24 * time.Hour) // 만료시간
	b := make([]byte, 16)
	rand.Read(b)                                  // 바이트가 랜덤하게 채워진다.
	state := base64.URLEncoding.EncodeToString(b) // 스트링으로 인코딩
	cookie := http.Cookie{
		Name: "oauthstate", Value: state, Expires: expiration,
	}
	http.SetCookie(w, &cookie)
	return state
}

// 콜백 요청이 왔을 때
func googleAuthCallback(w http.ResponseWriter, r *http.Request) {
	oauthstate, _ := r.Cookie("oauthstate")
	if r.FormValue("state") != oauthstate.Value { // 구글이 보낸 state값과 쿠키값이 다르면
		log.Printf("invalid google oauth state cookie: %s state: %s\n", oauthstate.Value, r.FormValue("state"))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect) // 루트 경로로 리다이렉트
		return
	}
	// 정상적인 요청이라면 구글측에서 보낸 코드를 통해 토큰을 만들어서 유저 정보를 가져옴
	data, err := getGoogleUserInfo(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprint(w, string(data)) // 유저에게 구글이 보낸 정보를 전송
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

// 구글이 보낸 코드를 통해 토큰 생성, oauthGoogleUrlAPI와 합쳐서 유저 정보 요청
func getGoogleUserInfo(code string) ([]byte, error) {
	// context는 쓰레드 간 데이터를 주고받을 때 사용되는 저장소
	// context.Background() == 기본 context - 멀티쓰레드 환경이 아니어서
	// Exchange의 반환값인 토큰에는 AccessToken, RefreshToken이 있다.
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Failed to Exchange %s", err.Error())
	}

	res, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get UserInfo %s", err.Error())
	}

	return ioutil.ReadAll(res.Body) // 유저 정보를 전송
}

func main() {
	mux := pat.New()
	mux.HandleFunc("/auth/google/login", googleLoginHandler)
	mux.HandleFunc("/auth/google/callback", googleAuthCallback)
	n := negroni.Classic()
	n.UseHandler(mux)

	http.ListenAndServe(":3000", n)
}
