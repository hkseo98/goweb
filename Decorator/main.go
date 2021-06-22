package main

// decorator 디자인 패턴!
// 각 컴포넌트에 다른 컴포넌트와, 실행 함수를 넣어서 가장 핵심기능 이외에 부가 기능들을 블록처럼 추가할 수 있게 하는 패턴

import (
	"fmt"

	"github.com/tuckersGo/goWeb/web9/cipher"
	"github.com/tuckersGo/goWeb/web9/lzw"
)

var sentData string
var recvData string

// 오퍼레이터 함수를 가지는 인터페이스
type Component interface {
	Operator(string) // 인자로 스트링을 받는 오퍼레이터 함수가 있음
}

// 기본 기능 컴포넌트
type SendComponent struct {
}

// 기본 데이터 전송 기능
func (self *SendComponent) Operator(data string) {
	// send data
	sentData = data
}

// 압축 기능 컴포넌트
type ZipComponent struct {
	com Component
}

// 압축 기능
func (self *ZipComponent) Operator(data string) {
	zipData, err := lzw.Write([]byte(data))
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(zipData))
}

// 암호화 기능 컴포넌트
type EncryptComponent struct {
	key string
	com Component
}

// 암호화 기능
func (self *EncryptComponent) Operator(data string) {
	endata, err := cipher.Encrypt([]byte(data), self.key)
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(endata))
}

// 복호와 기능 컴포넌트
type DecryptComponent struct {
	key string
	com Component
}

// 복호화 기능
func (self *DecryptComponent) Operator(data string) {
	deData, err := cipher.Decrypt([]byte(data), self.key)
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(deData))
}

// 압축 해제 기능 컴포넌트
type UnzipComponent struct {
	com Component
}

// 압축 해제 기능
func (self *UnzipComponent) Operator(data string) {
	unzipData, err := lzw.Read([]byte(data))
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(unzipData))
}

// 기본기능 컴포넌트
type ReadComponent struct{}

// 기본 읽기 기능
func (self *ReadComponent) Operator(data string) {
	recvData = data
}

func main() {
	// 순서를 바꾸거나 부가기능을 추가, 삭제 가능. 기본 기능은 변하지 않음
	sender := EncryptComponent{
		key: "abcde",
		com: &ZipComponent{
			com: &SendComponent{},
		},
	}
	sender.Operator("Hello World")
	fmt.Println(sentData)

	receiver := &UnzipComponent{
		com: &DecryptComponent{
			key: "abcde",
			com: &ReadComponent{},
		},
	}

	receiver.Operator(sentData)
	fmt.Println(recvData)
}
